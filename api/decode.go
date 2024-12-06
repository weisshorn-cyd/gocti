package api

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

var ErrNotAPointer = errors.New("result must be a pointer")

// NotImplementingError is returned when decoding a GraphQL interface into a type that does not implement it.
type NotImplementingError struct {
	InterfaceType      reflect.Type
	ImplementationType reflect.Type
}

func (e NotImplementingError) Error() string {
	return fmt.Sprintf(
		"'%s' does not implement the GraphQL interface '%s'",
		e.ImplementationType.Name(),
		e.InterfaceType.Name(),
	)
}

// UnimplementedDecodingError is returned when a decoding hook 'fromType -> toType' is missing.
type UnimplementedDecodingError struct {
	fromType reflect.Type
	toType   reflect.Type
	data     any
}

func (e UnimplementedDecodingError) Error() string {
	return fmt.Sprintf("missing decoding hook for '%v'(%v) -> '%v'", e.fromType, e.data, e.toType)
}

// Decode wraps mapstructure.Decode, using a custom DecoderConfig:
// It targets the "gocti" struct tag.
// It handles special types with a custom DecodeFunc.
func Decode(input, output any) error {
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Metadata:   nil,
			TagName:    "gocti",
			DecodeHook: goctiHookFunc,
			Result:     output,
		},
	)
	if err != nil {
		return fmt.Errorf("creating decoder: %w", err)
	}

	if err := decoder.Decode(input); err != nil {
		return fmt.Errorf("decoding: %w", err)
	}

	return nil
}

// goctiHookFunc is a mapstructure.DecodeHookFunc that handles the conversion of OpenCTI GraphQL types:
//   - DateTime (RFC3339 e.g. 2006-01-02T15:04:05Z07:00) -> time.Time
//
//nolint:wrapcheck // Additional information not required
func goctiHookFunc(fromType reflect.Type, toType reflect.Type, data any) (any, error) {
	// Handle DateTime fields resolving as time.Time
	if toType != reflect.TypeOf(time.Time{}) {
		return data, nil
	}

	//nolint:exhaustive // Errors out by default.
	switch fromType.Kind() {
	case reflect.String:
		//nolint:forcetypeassert // Reflect kind switch
		return time.Parse(time.RFC3339, data.(string))
	case reflect.Pointer:
		if fromType.Elem() == reflect.TypeFor[time.Time]() {
			//nolint:forcetypeassert // Reflect kind switch + type test
			return *data.(*time.Time), nil
		}

		fallthrough
	default:
		return data, UnimplementedDecodingError{
			fromType: fromType,
			toType:   toType,
			data:     data,
		}
	}
}

// DecodeInterface converts a GraphQL Interface type into one of its implementations.
// It uses the data returned by [GraphQLInterface.Remainder] to complete the fields
// specific to the implementation.
func DecodeInterface(input GraphQLInterface, output any) error {
	outputType := reflect.TypeOf(output)

	// Accept pointers only.
	if outputType.Kind() != reflect.Ptr {
		return ErrNotAPointer
	}

	outputType = outputType.Elem()

	// Check if output is part of implementation types.
	implements := false

	for _, impl := range input.Implementations() {
		if impl == outputType {
			implements = true
		}
	}

	if !implements {
		return NotImplementingError{reflect.TypeOf(input), outputType}
	}

	receivingMap := map[string]any{}

	err := Decode(input, &receivingMap)
	if err != nil {
		return fmt.Errorf("mapping interface: %w", err)
	}

	// Add the remaining fields to the map.
	for k, v := range input.Remainder() {
		receivingMap[k] = v
	}

	err = Decode(input, &output)
	if err != nil {
		return fmt.Errorf("decoding interface: %w", err)
	}

	return nil
}

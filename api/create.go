package api

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

type CreatableEntity interface {
	DefaultProperties() string
	CreateQueryString(customAttributes string) string
	CreateResponseField() string
}

type AddInput interface {
	Input() (map[string]any, error)
}

// Create will push the chosen [CreatableEntity] with the attributes given in 'input' to the OpenCTI server.
// If the entity already exists on the platform, it will be updated instead.
// It returns a map containing the chosen 'customAttributes' structure, as per GraphQL syntax.
func Create[E CreatableEntity](
	ctx context.Context,
	client Client,
	customAttributes string,
	input AddInput,
) (map[string]any, error) {
	// Retrieving entity-specific properties and query
	var entity E

	if customAttributes == "" {
		customAttributes = entity.DefaultProperties()
	}

	query := entity.CreateQueryString(customAttributes)

	// Querying OpenCTI server to create the entity
	client.Logger().Info("Creating entity",
		"entity", reflect.TypeFor[E]().Name(),
	)

	queryInput, err := input.Input()
	if err != nil {
		return nil, fmt.Errorf("failed to format input: %w", err)
	}

	queryData, err := client.Query(ctx, query, queryInput)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s: %w", reflect.TypeFor[E]().Name(), err)
	}

	// Processing the response: Expected struct: {"entity":map[string]any}
	resp, ok := queryData[entity.CreateResponseField()]
	if !ok {
		return nil, MissingFieldError{FieldName: entity.CreateResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

// StructuredCreate is a wrapper around [Create] that returns the data as a struct instead of a map.
// The provided ReturnStruct type must have a structure compatible with the GraphQL schema of the target entity.
// Field name aliases can (and should) be set using the "gocti" tag (e.g. struct{ OpenCTIName string `gocti:"name"` }).
// WARNING: Golang default values will be returned for fields that are not set or have a 'null' value on the server.
func StructuredCreate[E CreatableEntity, ReturnStruct any](
	ctx context.Context,
	client Client,
	customAttributes string,
	input AddInput,
) (ReturnStruct, error) {
	var (
		output     ReturnStruct
		outputType = reflect.TypeFor[ReturnStruct]()
	)

	// Accept structs only.
	if outputType.Kind() != reflect.Struct {
		return output, fmt.Errorf("output type error: %w", InterfaceTypeError{reflect.Struct, outputType})
	}

	// Retrieving the entity from the server.
	data, err := Create[E](ctx, client, customAttributes, input)
	if err != nil {
		return output, err
	}

	// Decoding the returned data into the target struct.
	if err := Decode(data, &output); err != nil {
		return output, fmt.Errorf("decoding output struct: %w", err)
	}

	return output, nil
}

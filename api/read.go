package api

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

type ReadableEntity interface {
	DefaultProperties() string
	ReadQueryString(customAttributes string) string
	ReadResponseField() string
}

// Read retrieves the chosen [ReadableEntity] with the provided 'id' from the OpenCTI server.
// It returns a map containing the chosen 'customAttributes' structure, as per GraphQL syntax.
func Read[E ReadableEntity](
	ctx context.Context,
	client Client,
	customAttributes, id string,
) (map[string]any, error) {
	// Retrieving entity-specific properties and query
	var entity E

	if customAttributes == "" {
		customAttributes = entity.DefaultProperties()
	}

	query := entity.ReadQueryString(customAttributes)

	// Querying OpenCTI server for list of entities
	client.Logger().Info("Reading entity",
		"entity", reflect.TypeFor[E]().Name(),
		"id", id,
	)

	queryData, err := client.Query(ctx, query, map[string]any{"id": id})
	if err != nil {
		return nil, fmt.Errorf("unable to read entity: %w", err)
	}

	// Processing the response: Expected struct: {"entity":map[string]any}
	resp, ok := queryData[entity.ReadResponseField()]
	if !ok {
		return nil, MissingFieldError{FieldName: entity.ReadResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

// StructuredRead is a wrapper around [Read] that returns the data as a struct instead of a map.
// The provided ReturnStruct type must have a structure compatible with the GraphQL schema of the target entity.
// Field name aliases can (and should) be set using the "gocti" tag (e.g. struct{ OpenCTIName string `gocti:"name"` }).
// WARNING: Golang default values will be returned for fields that are not set or have a 'null' value on the server.
func StructuredRead[E ReadableEntity, ReturnStruct any](
	ctx context.Context,
	client Client,
	customAttributes, id string,
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
	data, err := Read[E](ctx, client, customAttributes, id)
	if err != nil {
		return output, err
	}

	// Decoding the returned data into the target struct.
	if err := Decode(data, &output); err != nil {
		return output, fmt.Errorf("decoding output struct: %w", err)
	}

	return output, nil
}

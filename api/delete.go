package api

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

type DeletableEntity interface {
	DeleteQueryString() string
	DeleteResponseField() string
}

// Delete deletes the chosen DeletableEntity with the provided 'id' from the OpenCTI server.
// It returns the ID of the recently deleted entity.
func Delete[E DeletableEntity](
	ctx context.Context,
	client Client,
	id string,
) (string, error) {
	// Retrieving entity-specific properties and query
	var entity E

	query := entity.DeleteQueryString()

	// Querying OpenCTI server for deleting the entity
	client.Logger().Info("Deleting entity",
		"entity", reflect.TypeFor[E]().Name(),
		"id", id,
	)

	queryData, err := client.Query(ctx, query, map[string]any{"id": id})
	if err != nil {
		return "", fmt.Errorf("unable to delete entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":{"delete":string}} or {"query_name":string}
	resp, ok := queryData[entity.DeleteResponseField()]
	if !ok {
		return "", MissingFieldError{FieldName: entity.DeleteResponseField()}
	}

	// Testing for the first possible return structure.
	returnedMap := struct {
		Delete string `mapstructure:"delete"`
	}{}
	if err := mapstructure.Decode(resp, &returnedMap); err == nil {
		return returnedMap.Delete, nil
	}

	// If it fails, testing for the second possible return structure.
	returnedID, ok := resp.(string)
	if !ok {
		return "", TypeAssertionError{resp, "string"}
	}

	return returnedID, nil
}

package api

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"

	"github.com/weisshorn-cyd/gocti/graphql"
	"github.com/weisshorn-cyd/gocti/list"
)

type ListableEntity interface {
	DefaultProperties() string
	ListQueryString(customAttributes string) string
	ListResponseField() string
}

type listResponse struct {
	Edges    []node           `mapstructure:"edges"`
	PageInfo graphql.PageInfo `mapstructure:"pageInfo"`
}

type node struct {
	Node map[string]any `mapstructure:"node"`
}

// List fetches a list of the chosen [ListableEntity] from the OpenCTI server
// It returns a slice of maps, each containing the chosen 'customAttributes' structure, as per GraphQL syntax.
// Pagination information can be retrieved by providing a non-nil pointer to a [graphql.PageInfo].
func List[E ListableEntity](
	ctx context.Context,
	client Client,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]map[string]any, error) {
	if getAll {
		return listAll[E](ctx, client, customAttributes, pageInfo, opts...)
	}

	// Retrieving entity-specific parameters
	query, responseField, vars, err := resolveParameters[E](client, customAttributes, opts...)
	if err != nil {
		return nil, fmt.Errorf("resolving parameters: %w", err)
	}

	// Querying the GraphQL server
	queryData, err := client.Query(ctx, query, vars)
	if err != nil {
		return nil, fmt.Errorf("unable to list entities: %w", err)
	}

	// Processing the response: Expected struct: {"entities":listResponseField}
	resp, ok := queryData[responseField]
	if !ok {
		return nil, MissingFieldError{FieldName: responseField}
	}

	listResp := listResponse{}
	if err := mapstructure.Decode(resp, &listResp); err != nil {
		return nil, fmt.Errorf("failed to retrieve 'listResponse' struct: %w", err)
	}

	finalList := []map[string]any{}
	for _, node := range listResp.Edges {
		finalList = append(finalList, node.Node)
	}

	if pageInfo != nil {
		*pageInfo = listResp.PageInfo
	}

	return finalList, nil
}

// StructuredList is a wrapper around [List] that returns the data as a list of structs instead of a list of maps.
// The provided ReturnStruct type must have a structure compatible with the GraphQL schema of the target entity.
// Field name aliases can (and should) be set using the "gocti" tag (e.g. struct{ OpenCTIName string `gocti:"name"` }).
// WARNING: Golang default values will be returned for fields that are not set or have a 'null' value on the server.
func StructuredList[E ListableEntity, ReturnStruct any](
	ctx context.Context,
	client Client,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]ReturnStruct, error) {
	returnType := reflect.TypeFor[ReturnStruct]()

	// Accept structs only.
	if returnType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("output type error: %w", InterfaceTypeError{reflect.Struct, returnType})
	}

	// Retrieving the list of entities from the server.
	data, err := List[E](ctx, client, customAttributes, getAll, pageInfo, opts...)
	if err != nil {
		return nil, err
	}

	// Decoding the returned data into the target struct.
	output := make([]ReturnStruct, len(data))
	for index, item := range data {
		if err := Decode(item, &output[index]); err != nil {
			return nil, fmt.Errorf("decoding output struct: %w", err)
		}
	}

	return output, nil
}

// listAll will repetitively call [List] to retrieve the complete list of entities from the server.
func listAll[E ListableEntity](
	ctx context.Context,
	client Client,
	customAttributes string,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]map[string]any, error) {
	pagination := graphql.PageInfo{HasNextPage: true}
	entityList := []map[string]any{}

	count := 0

	for pagination.HasNextPage {
		// Move cursor
		opts = append(opts, list.WithAfter(pagination.EndCursor))

		// Query
		tempList, err := List[E](ctx, client, customAttributes, false, &pagination, opts...)
		if err != nil {
			return nil, err
		}

		// Fill in pageInfo if provided
		if pageInfo != nil && pageInfo.StartCursor == "" {
			pageInfo.StartCursor = pagination.StartCursor
			pageInfo.HasPreviousPage = pagination.HasPreviousPage
		}

		count += len(tempList)
		client.Logger().Info("Listing all entities",
			"count", count,
			"total", pagination.GlobalCount,
		)

		entityList = append(entityList, tempList...)
	}

	// Update returned pageInfo to encompass all data
	if pageInfo != nil {
		pageInfo.GlobalCount = pagination.GlobalCount
		pageInfo.EndCursor = pagination.EndCursor
		pageInfo.HasNextPage = false
	}

	return entityList, nil
}

// resolveParameters will return the custom attributes and options to be used for a list operation.
// Return values: List query string (with custom attributes), GraphQL response field name, configured variables struct.
func resolveParameters[E ListableEntity](
	defaults Defaults,
	customAttributes string,
	opts ...list.Option,
) (string, string, map[string]any, error) {
	var entity E

	// Resolving query string
	if customAttributes == "" {
		customAttributes = entity.DefaultProperties()
	}

	query := entity.ListQueryString(customAttributes)

	// Retrieving GraphQL response field name
	responseField := entity.ListResponseField()

	// Gather all config Options from the Defaults Configuration
	configOptions := []list.Option{}

	if pageSize, ok := defaults.DefaultPageSize(); ok {
		configOptions = append(configOptions, list.WithFirst(pageSize))
	}

	if orderBy, ok := defaults.DefaultOrderBy(); ok {
		configOptions = append(configOptions, list.WithOrderBy(orderBy))
	}

	if orderMode, ok := defaults.DefaultOrderMode(); ok {
		configOptions = append(configOptions, list.WithOrderMode(list.OrderingMode(orderMode)))
	}

	// Prepend all default config Options to user defined Options
	configOptions = append(configOptions, opts...)

	// Apply all options to a [list.QueryVars] struct
	queryVars := list.NewQueryVars()
	for _, opt := range configOptions {
		opt(queryVars)
	}

	// Converting query vars to map
	vars, err := queryVars.Mapping()
	if err != nil {
		return "", "", nil, fmt.Errorf("converting query vars to map: %w", err)
	}

	return query, responseField, vars, nil
}

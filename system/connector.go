package system

import (
	"context"
	"errors"
	"fmt"

	"github.com/weisshorn-cyd/gocti/api"

	_ "embed"
)

//go:embed misc_queries/connector_list.graphql
var connectorsQueryString string

// ConnectorWithConfig represents the GraphQL input type ConnectorWithConfig,
// used by uploadAndAskJobImport to select which connector(s) process an uploaded file.
type ConnectorWithConfig struct {
	ConnectorID   string `json:"connectorId"`
	Configuration string `json:"configuration,omitempty"`
}

var ErrConnectorNotFound = errors.New("connector not found")

type Connector struct {
	ID     string `gocti:"id"`
	Name   string `gocti:"name"`
	Active bool   `gocti:"active"`
}

// GetConnectorByName lists connectors and returns the first whose name matches.
func GetConnectorByName(
	ctx context.Context,
	client api.Client,
	name string,
) (Connector, error) {
	data, err := client.Query(ctx, connectorsQueryString, nil)
	if err != nil {
		return Connector{}, fmt.Errorf("cannot list connectors: %w", err)
	}

	resp, ok := data["connectors"]
	if !ok {
		return Connector{}, api.MissingFieldError{FieldName: "connectors"}
	}

	connectors := []Connector{}
	if err := api.Decode(resp, &connectors); err != nil {
		return Connector{}, fmt.Errorf("failed to decode connectors: %w", err)
	}

	for _, c := range connectors {
		if c.Name == name {
			return c, nil
		}
	}

	return Connector{}, fmt.Errorf("%w: %s", ErrConnectorNotFound, name)
}

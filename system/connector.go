package system

import (
	"context"
	"errors"
	"fmt"

	"github.com/weisshorn-cyd/gocti/api"

	_ "embed"
)

//go:embed list_queries/connector_list_query.txt
var connectorsQueryString string

// ConnectorWithConfig represents the GraphQL input type ConnectorWithConfig,
// used by uploadAndAskJobImport to select which connector(s) process an uploaded file.
type ConnectorWithConfig struct {
	ConnectorID   string `json:"connectorId"`
	Configuration string `json:"configuration,omitempty"`
}

var ErrConnectorNotFound = errors.New("connector not found")

type connector struct {
	ID   string `gocti:"id"`
	Name string `gocti:"name"`
}

// GetConnectorByName lists connectors and returns the first whose name matches.
func GetConnectorByName(
	ctx context.Context,
	client api.Client,
	name string,
) (string, error) {
	data, err := client.Query(ctx, connectorsQueryString, nil)
	if err != nil {
		return "", fmt.Errorf("cannot list connectors: %w", err)
	}

	resp, ok := data["connectors"]
	if !ok {
		return "", api.MissingFieldError{FieldName: "connectors"}
	}

	connectors := []connector{}
	if err := api.Decode(resp, &connectors); err != nil {
		return "", fmt.Errorf("failed to decode connectors: %w", err)
	}

	for _, c := range connectors {
		if c.Name == name {
			return c.ID, nil
		}
	}

	return "", fmt.Errorf("%w: %s", ErrConnectorNotFound, name)
}

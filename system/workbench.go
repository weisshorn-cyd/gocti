package system

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/graphql"

	_ "embed"
)

type Workbench struct {
	ID           string        `gocti:"id"`
	Name         string        `gocti:"name"`
	Size         int           `gocti:"size"`
	LastModified time.Time     `gocti:"lastModified"`
	UploadStatus graphql.State `gocti:"uploadStatus"`
	URL          string        `gocti:"-"`
}

//go:embed misc_queries/uploadPending.graphql
var uploadPendingQueryString string

// CreateWorkbench creates a new Workbench on the server based on the provided file.
// The [api.File] struct must contain all Name, Data, and MIME-type.
// Workbench name will match the file name.
// File Data is typically a valid stix bundle formatted as JSON.
func CreateWorkbench(
	ctx context.Context,
	client api.Client,
	file api.File,
	errOnExisting bool,
) (Workbench, error) {
	variables := map[string]any{
		"file":            file,
		"errorOnExisting": errOnExisting,
	}

	data, err := client.Query(ctx, uploadPendingQueryString, variables)
	if err != nil {
		return Workbench{}, fmt.Errorf("cannot create workbench: %w", err)
	}

	resp, ok := data["uploadPending"]
	if !ok {
		return Workbench{}, api.MissingFieldError{FieldName: "uploadPending"}
	}

	output := Workbench{}
	if err := api.Decode(resp, &output); err != nil {
		return output, fmt.Errorf("failed to retrieve result map: %w", err)
	}

	url := base64.StdEncoding.EncodeToString([]byte(output.ID))
	output.URL = "/dashboard/data/import/pending/" + strings.TrimSuffix(url, "==")

	return output, nil
}

//go:embed misc_queries/deleteImport.graphql
var deleteImportQueryString string

// DeleteWorkbench removes an active workbench from the server.
// id is typically a path, as returned by [CreateWorkbench] e.g. "import/pending/ExempleWorkbench.json".
func DeleteWorkbench(
	ctx context.Context,
	client api.Client,
	id string,
) (string, error) {
	data, err := client.Query(ctx, deleteImportQueryString, map[string]any{
		"fileName": id,
	})
	if err != nil {
		return "", fmt.Errorf("failed to delete workbench: %w", err)
	}

	resp, ok := data["deleteImport"]
	if !ok {
		return "", api.MissingFieldError{FieldName: "uploadPending"}
	}

	if id, ok = resp.(string); !ok {
		return "", fmt.Errorf("failed to retrieve id: %w", err)
	}

	return id, nil
}

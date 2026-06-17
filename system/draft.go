package system

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/graphql"

	_ "embed"
)

type Draft struct {
	ID          string `gocti:"id"`
	Name        string `gocti:"name"`
	DraftStatus string `gocti:"draft_status"`
	CreatedAt   string `gocti:"created_at"`
	URL         string `gocti:"url"`
}

type DraftWorkspaceAddInput struct {
	Name              string              `json:"name"`
	Description       *string             `json:"description,omitempty"`
	EntityID          *string             `json:"entity_id,omitempty"`
	ObjectAssignee    []string            `json:"objectAssignee,omitempty"`
	ObjectParticipant []string            `json:"objectParticipant,omitempty"`
	CreatedBy         *string             `json:"createdBy,omitempty"`
	AuthorizedMembers []MemberAccessInput `json:"authorized_members,omitempty"`
}

type MemberAccessInput struct {
	ID                   string   `json:"id"`
	AccessRight          string   `json:"access_right"`
	GroupsRestrictionIDs []string `json:"groups_restriction_ids,omitempty"`
}

//go:embed misc_queries/createDraft.graphql
var createDraftQueryString string

// CreateDraft creates a new Draft on the server.
func CreateDraft(
	ctx context.Context,
	client api.Client,
	input DraftWorkspaceAddInput,
) (Draft, error) {
	variables := map[string]any{
		"input": input,
	}

	data, err := client.Query(ctx, createDraftQueryString, variables)
	if err != nil {
		return Draft{}, fmt.Errorf("cannot create draft: %w", err)
	}

	resp, ok := data["draftWorkspaceAdd"]
	if !ok {
		return Draft{}, api.MissingFieldError{FieldName: "draftWorkspaceAdd"}
	}

	output := Draft{}
	if err := api.Decode(resp, &output); err != nil {
		return output, fmt.Errorf("failed to retrieve result map: %w", err)
	}

	url := base64.StdEncoding.EncodeToString([]byte(output.ID))
	output.URL = "/dashboard/data/import/draft/" + strings.TrimSuffix(url, "==")

	return output, nil
}

//go:embed misc_queries/deleteDraft.graphql
var deleteDraftQueryString string

// DeleteDraft removes a draft from the server.
func DeleteDraft(
	ctx context.Context,
	client api.Client,
	id string,
) (string, error) {
	data, err := client.Query(ctx, deleteDraftQueryString, map[string]any{
		"id": id,
	})
	if err != nil {
		return "", fmt.Errorf("failed to delete draft: %w", err)
	}

	resp, ok := data["draftWorkspaceDelete"]
	if !ok {
		return "", api.MissingFieldError{FieldName: "draftWorkspaceDelete"}
	}

	if id, ok = resp.(string); !ok {
		return "", fmt.Errorf("failed to retrieve id: %w", err)
	}

	return id, nil
}

//go:embed misc_queries/importFile.graphql
var importFileQueryString string

// ImportFile uploads a file into an existing draft on the server.
// The [api.File] struct must contain all Name, Data, and MIME-type.
// File Data is typically a valid stix bundle formatted as JSON.
func ImportFile(
	ctx context.Context,
	client api.Client,
	draftID string,
	file api.File,
	connectors []ConnectorWithConfig,
) (graphql.File, error) {
	variables := map[string]any{
		"file":            file,
		"fileMarkings":    []string{},
		"connectors":      connectors,
		"validationMode":  "draft",
		"draftId":         draftID,
		"noTriggerImport": true,
	}

	data, err := client.Query(ctx, importFileQueryString, variables)
	if err != nil {
		return graphql.File{}, fmt.Errorf("cannot upload file: %w", err)
	}

	resp, ok := data["uploadAndAskJobImport"]
	if !ok {
		return graphql.File{}, api.MissingFieldError{FieldName: "uploadAndAskJobImport"}
	}

	output := graphql.File{}
	if err := api.Decode(resp, &output); err != nil {
		return output, fmt.Errorf("failed to retrieve result map: %w", err)
	}

	return output, nil
}

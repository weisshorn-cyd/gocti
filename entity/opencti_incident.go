// Code generated by '/tools/gocti_type_generator' for OpenCTI version 6.7.8 - DO NOT EDIT.

package entity

import (
	"fmt"
	"time"

	"github.com/weisshorn-cyd/gocti/graphql"

	_ "embed"
)

type Incident struct {
	graphql.Incident `gocti:",squash"`
}

//go:embed default_properties/incident_default_properties.txt
var incidentDefaultProperties string

func (i Incident) DefaultProperties() string {
	return incidentDefaultProperties
}

// Implementing the [api.ListableEntity] interface.

//go:embed list_queries/incident_list_query.txt
var incidentListQueryString string

func (i Incident) ListQueryString(customAttributes string) string {
	return fmt.Sprintf(
		incidentListQueryString,
		customAttributes,
	)
}

func (i Incident) ListResponseField() string { return "incidents" }

// Implementing the [api.ReadableEntity] interface.

func (i Incident) ReadQueryString(customAttributes string) string {
	return fmt.Sprintf(
		`query ($id: String) {
            incident (id: $id) {%s}
        }`,
		customAttributes,
	)
}

func (i Incident) ReadResponseField() string { return "incident" }

// Implementing the [api.CreatableEntity] interface.

func (i Incident) CreateQueryString(customAttributes string) string {
	return fmt.Sprintf(
		`mutation ($input: IncidentAddInput!) {
            incidentAdd (input: $input) {%s}
        }`,
		customAttributes,
	)
}

func (i Incident) CreateResponseField() string { return "incidentAdd" }

// IncidentAddInput represents a GraphQL INPUT_OBJECT
// Some fields from the OpenCTI schema may be missing
// (See the examples for ways to expand an existing type).
type IncidentAddInput struct {
	// IncidentAddInput
	StixID             string     `gocti:"stix_id"               json:"stix_id,omitempty"`
	XOpenctiStixIDs    []string   `gocti:"x_opencti_stix_ids"    json:"x_opencti_stix_ids,omitempty"`
	Name               string     `gocti:"name"                  json:"name,omitempty"`
	Description        string     `gocti:"description"           json:"description,omitempty"`
	Confidence         int        `gocti:"confidence"            json:"confidence,omitempty"`
	Revoked            bool       `gocti:"revoked"               json:"revoked,omitempty"`
	Lang               string     `gocti:"lang"                  json:"lang,omitempty"`
	Objective          string     `gocti:"objective"             json:"objective,omitempty"`
	FirstSeen          *time.Time `gocti:"first_seen"            json:"first_seen,omitempty"`
	LastSeen           *time.Time `gocti:"last_seen"             json:"last_seen,omitempty"`
	Aliases            []string   `gocti:"aliases"               json:"aliases,omitempty"`
	IncidentType       string     `gocti:"incident_type"         json:"incident_type,omitempty"`
	Severity           string     `gocti:"severity"              json:"severity,omitempty"`
	Source             string     `gocti:"source"                json:"source,omitempty"`
	CreatedBy          string     `gocti:"createdBy"             json:"createdBy,omitempty"`
	ObjectOrganization []string   `gocti:"objectOrganization"    json:"objectOrganization,omitempty"`
	ObjectMarking      []string   `gocti:"objectMarking"         json:"objectMarking,omitempty"`
	ObjectAssignee     []string   `gocti:"objectAssignee"        json:"objectAssignee,omitempty"`
	ObjectParticipant  []string   `gocti:"objectParticipant"     json:"objectParticipant,omitempty"`
	ObjectLabel        []string   `gocti:"objectLabel"           json:"objectLabel,omitempty"`
	ExternalReferences []string   `gocti:"externalReferences"    json:"externalReferences,omitempty"`
	Created            *time.Time `gocti:"created"               json:"created,omitempty"`
	Modified           *time.Time `gocti:"modified"              json:"modified,omitempty"`
	XOpenctiWorkflowID string     `gocti:"x_opencti_workflow_id" json:"x_opencti_workflow_id,omitempty"`
	ClientMutationID   string     `gocti:"clientMutationId"      json:"clientMutationId,omitempty"`
	Update             bool       `gocti:"update"                json:"update,omitempty"`
	File               []byte     `gocti:"file"                  json:"file,omitempty"`
}

func (input IncidentAddInput) Input() (map[string]any, error) {
	return map[string]any{
		"input": input,
	}, nil
}

// Implementing the [api.DeletableEntity] interface.

func (i Incident) DeleteQueryString() string {
	return `mutation ($id: ID!) {
                incidentEdit (id: $id) {
                    delete
                }
            }`
}

func (i Incident) DeleteResponseField() string { return "incidentEdit" }

// Code generated by '/tools/gocti_type_generator' for OpenCTI version 6.7.8 - DO NOT EDIT.

package entity

import (
	"fmt"
	"time"

	"github.com/weisshorn-cyd/gocti/graphql"

	_ "embed"
)

type CaseRfi struct {
	graphql.CaseRfi `gocti:",squash"`
}

//go:embed default_properties/case_rfi_default_properties.txt
var caseRfiDefaultProperties string

func (c CaseRfi) DefaultProperties() string {
	return caseRfiDefaultProperties
}

// Implementing the [api.ListableEntity] interface.

//go:embed list_queries/case_rfi_list_query.txt
var caseRfiListQueryString string

func (c CaseRfi) ListQueryString(customAttributes string) string {
	return fmt.Sprintf(
		caseRfiListQueryString,
		customAttributes,
	)
}

func (c CaseRfi) ListResponseField() string { return "caseRfis" }

// Implementing the [api.ReadableEntity] interface.

func (c CaseRfi) ReadQueryString(customAttributes string) string {
	return fmt.Sprintf(
		`query ($id: String!) {
            caseRfi (id: $id) {%s}
        }`,
		customAttributes,
	)
}

func (c CaseRfi) ReadResponseField() string { return "caseRfi" }

// Implementing the [api.CreatableEntity] interface.

func (c CaseRfi) CreateQueryString(customAttributes string) string {
	return fmt.Sprintf(
		`mutation ($input: CaseRfiAddInput!) {
            caseRfiAdd (input: $input) {%s}
        }`,
		customAttributes,
	)
}

func (c CaseRfi) CreateResponseField() string { return "caseRfiAdd" }

// CaseRfiAddInput represents a GraphQL INPUT_OBJECT
// Some fields from the OpenCTI schema may be missing
// (See the examples for ways to expand an existing type).
type CaseRfiAddInput struct {
	// CaseRfiAddInput
	StixID                string                      `gocti:"stix_id"                  json:"stix_id,omitempty"`
	XOpenctiStixIDs       []string                    `gocti:"x_opencti_stix_ids"       json:"x_opencti_stix_ids,omitempty"`
	Name                  string                      `gocti:"name"                     json:"name,omitempty"`
	Description           string                      `gocti:"description"              json:"description,omitempty"`
	Content               string                      `gocti:"content"                  json:"content,omitempty"`
	ContentMapping        string                      `gocti:"content_mapping"          json:"content_mapping,omitempty"`
	Severity              string                      `gocti:"severity"                 json:"severity,omitempty"`
	Priority              string                      `gocti:"priority"                 json:"priority,omitempty"`
	Confidence            int                         `gocti:"confidence"               json:"confidence,omitempty"`
	Revoked               bool                        `gocti:"revoked"                  json:"revoked,omitempty"`
	Lang                  string                      `gocti:"lang"                     json:"lang,omitempty"`
	Objects               []string                    `gocti:"objects"                  json:"objects,omitempty"`
	CreatedBy             string                      `gocti:"createdBy"                json:"createdBy,omitempty"`
	ObjectMarking         []string                    `gocti:"objectMarking"            json:"objectMarking,omitempty"`
	ObjectAssignee        []string                    `gocti:"objectAssignee"           json:"objectAssignee,omitempty"`
	ObjectParticipant     []string                    `gocti:"objectParticipant"        json:"objectParticipant,omitempty"`
	ObjectOrganization    []string                    `gocti:"objectOrganization"       json:"objectOrganization,omitempty"`
	ObjectLabel           []string                    `gocti:"objectLabel"              json:"objectLabel,omitempty"`
	ExternalReferences    []string                    `gocti:"externalReferences"       json:"externalReferences,omitempty"`
	Created               *time.Time                  `gocti:"created"                  json:"created,omitempty"`
	Modified              *time.Time                  `gocti:"modified"                 json:"modified,omitempty"`
	XOpenctiWorkflowID    string                      `gocti:"x_opencti_workflow_id"    json:"x_opencti_workflow_id,omitempty"`
	File                  []byte                      `gocti:"file"                     json:"file,omitempty"`
	ClientMutationID      string                      `gocti:"clientMutationId"         json:"clientMutationId,omitempty"`
	Update                bool                        `gocti:"update"                   json:"update,omitempty"`
	InformationTypes      []string                    `gocti:"information_types"        json:"information_types,omitempty"`
	CaseTemplates         []string                    `gocti:"caseTemplates"            json:"caseTemplates,omitempty"`
	AuthorizedMembers     []graphql.MemberAccessInput `gocti:"authorized_members"       json:"authorized_members,omitempty"`
	XOpenctiRequestAccess string                      `gocti:"x_opencti_request_access" json:"x_opencti_request_access,omitempty"`
}

func (input CaseRfiAddInput) Input() (map[string]any, error) {
	return map[string]any{
		"input": input,
	}, nil
}

// Implementing the [api.DeletableEntity] interface.

func (c CaseRfi) DeleteQueryString() string {
	return `mutation ($id: ID!) {
                caseRfiDelete (id: $id)
            }`
}

func (c CaseRfi) DeleteResponseField() string { return "caseRfiDelete" }

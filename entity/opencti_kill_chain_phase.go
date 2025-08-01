// Code generated by '/tools/gocti_type_generator' for OpenCTI version 6.7.8 - DO NOT EDIT.

package entity

import (
	"fmt"
	"time"

	"github.com/weisshorn-cyd/gocti/graphql"

	_ "embed"
)

type KillChainPhase struct {
	graphql.KillChainPhase `gocti:",squash"`
}

//go:embed default_properties/kill_chain_phase_default_properties.txt
var killChainPhaseDefaultProperties string

func (k KillChainPhase) DefaultProperties() string {
	return killChainPhaseDefaultProperties
}

// Implementing the [api.ListableEntity] interface.

//go:embed list_queries/kill_chain_phase_list_query.txt
var killChainPhaseListQueryString string

func (k KillChainPhase) ListQueryString(customAttributes string) string {
	return fmt.Sprintf(
		killChainPhaseListQueryString,
		customAttributes,
	)
}

func (k KillChainPhase) ListResponseField() string { return "killChainPhases" }

// Implementing the [api.ReadableEntity] interface.

func (k KillChainPhase) ReadQueryString(customAttributes string) string {
	return fmt.Sprintf(
		`query ($id: String!) {
            killChainPhase (id: $id) {%s}
        }`,
		customAttributes,
	)
}

func (k KillChainPhase) ReadResponseField() string { return "killChainPhase" }

// Implementing the [api.CreatableEntity] interface.

func (k KillChainPhase) CreateQueryString(customAttributes string) string {
	return fmt.Sprintf(
		`mutation ($input: KillChainPhaseAddInput!) {
            killChainPhaseAdd (input: $input) {%s}
        }`,
		customAttributes,
	)
}

func (k KillChainPhase) CreateResponseField() string { return "killChainPhaseAdd" }

// KillChainPhaseAddInput represents a GraphQL INPUT_OBJECT
// Some fields from the OpenCTI schema may be missing
// (See the examples for ways to expand an existing type).
type KillChainPhaseAddInput struct {
	// KillChainPhaseAddInput
	StixID           string     `gocti:"stix_id"            json:"stix_id,omitempty"`
	XOpenctiStixIDs  []string   `gocti:"x_opencti_stix_ids" json:"x_opencti_stix_ids,omitempty"`
	KillChainName    string     `gocti:"kill_chain_name"    json:"kill_chain_name,omitempty"`
	PhaseName        string     `gocti:"phase_name"         json:"phase_name,omitempty"`
	XOpenctiOrder    int        `gocti:"x_opencti_order"    json:"x_opencti_order,omitempty"`
	Created          *time.Time `gocti:"created"            json:"created,omitempty"`
	Modified         *time.Time `gocti:"modified"           json:"modified,omitempty"`
	ClientMutationID string     `gocti:"clientMutationId"   json:"clientMutationId,omitempty"`
	Update           bool       `gocti:"update"             json:"update,omitempty"`
}

func (input KillChainPhaseAddInput) Input() (map[string]any, error) {
	return map[string]any{
		"input": input,
	}, nil
}

// Implementing the [api.DeletableEntity] interface.

func (k KillChainPhase) DeleteQueryString() string {
	return `mutation ($id: ID!) {
                killChainPhaseEdit (id: $id) {
                    delete
                }
            }`
}

func (k KillChainPhase) DeleteResponseField() string { return "killChainPhaseEdit" }

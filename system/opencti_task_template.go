// Code generated by '/tools/gocti_type_generator' for OpenCTI version 6.7.8 - DO NOT EDIT.

package system

import (
	"fmt"

	"github.com/weisshorn-cyd/gocti/graphql"

	_ "embed"
)

type TaskTemplate struct {
	graphql.TaskTemplate `gocti:",squash"`
}

//go:embed default_properties/task_template_default_properties.txt
var taskTemplateDefaultProperties string

func (t TaskTemplate) DefaultProperties() string {
	return taskTemplateDefaultProperties
}

// Implementing the [api.ListableEntity] interface.

//go:embed list_queries/task_template_list_query.txt
var taskTemplateListQueryString string

func (t TaskTemplate) ListQueryString(customAttributes string) string {
	return fmt.Sprintf(
		taskTemplateListQueryString,
		customAttributes,
	)
}

func (t TaskTemplate) ListResponseField() string { return "taskTemplates" }

// Implementing the [api.ReadableEntity] interface.

func (t TaskTemplate) ReadQueryString(customAttributes string) string {
	return fmt.Sprintf(
		`query ($id: String!) {
            taskTemplate (id: $id) {%s}
        }`,
		customAttributes,
	)
}

func (t TaskTemplate) ReadResponseField() string { return "taskTemplate" }

// Implementing the [api.CreatableEntity] interface.

func (t TaskTemplate) CreateQueryString(customAttributes string) string {
	return fmt.Sprintf(
		`mutation ($input: TaskTemplateAddInput!) {
            taskTemplateAdd (input: $input) {%s}
        }`,
		customAttributes,
	)
}

func (t TaskTemplate) CreateResponseField() string { return "taskTemplateAdd" }

// TaskTemplateAddInput represents a GraphQL INPUT_OBJECT
// Some fields from the OpenCTI schema may be missing
// (See the examples for ways to expand an existing type).
type TaskTemplateAddInput struct {
	// TaskTemplateAddInput
	Name        string `gocti:"name"        json:"name,omitempty"`
	Description string `gocti:"description" json:"description,omitempty"`
}

func (input TaskTemplateAddInput) Input() (map[string]any, error) {
	return map[string]any{
		"input": input,
	}, nil
}

// Implementing the [api.DeletableEntity] interface.

func (t TaskTemplate) DeleteQueryString() string {
	return `mutation ($id: ID!) {
                taskTemplateDelete (id: $id)
            }`
}

func (t TaskTemplate) DeleteResponseField() string { return "taskTemplateDelete" }

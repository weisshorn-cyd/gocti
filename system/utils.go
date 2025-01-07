package system

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"

	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/entity"

	_ "embed"
)

// User utils

func (u User) EditResponseField() string { return "userEdit" }

//go:embed edit_queries/user_assign_group.graphql
var userAssignGroupQueryString string

func (u User) AssignGroup(
	ctx context.Context,
	client api.Client,
	groupID string,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id": u.ID,
		"input": entity.StixCoreRelationshipAddInput{
			RelationshipType: "member-of",
			ToID:             groupID,
		},
	}

	queryData, err := client.Query(
		ctx,
		userAssignGroupQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[u.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: u.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

//go:embed edit_queries/user_unassign_group.graphql
var userUnassignGroupQueryString string

func (u User) UnassignGroup(
	ctx context.Context,
	client api.Client,
	groupID string,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id":                u.ID,
		"relationship_type": "member-of",
		"toId":              groupID,
	}

	queryData, err := client.Query(
		ctx,
		userUnassignGroupQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[u.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: u.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

// Role utils

func (r Role) EditResponseField() string { return "roleEdit" }

//go:embed edit_queries/role_assign_capability.graphql
var roleAssignCapabilityQueryString string

func (r Role) AssignCapability(
	ctx context.Context,
	client api.Client,
	capabilityID string,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id": r.ID,
		"input": entity.StixCoreRelationshipAddInput{
			RelationshipType: "has-capability",
			ToID:             capabilityID,
		},
	}

	queryData, err := client.Query(
		ctx,
		roleAssignCapabilityQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[r.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: r.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

//go:embed edit_queries/role_unassign_capability.graphql
var roleUnassignCapabilityQueryString string

func (r Role) UnassignCapability(
	ctx context.Context,
	client api.Client,
	capabilityID string,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id":                r.ID,
		"relationship_type": "has-capability",
		"toId":              capabilityID,
	}

	queryData, err := client.Query(
		ctx,
		roleUnassignCapabilityQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[r.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: r.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

// Group utils

func (g Group) EditResponseField() string { return "groupEdit" }

//go:embed edit_queries/group_assign_marking_definition.graphql
var groupAssignMarkingDefinitionQueryString string

func (g Group) AssignMarkingDefinition(
	ctx context.Context,
	client api.Client,
	markingDefinitionID string,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id": g.ID,
		"input": entity.StixCoreRelationshipAddInput{
			RelationshipType: "accesses-to",
			ToID:             markingDefinitionID,
		},
	}

	queryData, err := client.Query(
		ctx,
		groupAssignMarkingDefinitionQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[g.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: g.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

//go:embed edit_queries/group_unassign_marking_definition.graphql
var groupUnassignMarkingDefinitionQueryString string

func (g Group) UnassignMarkingDefinition(
	ctx context.Context,
	client api.Client,
	markingDefinitionID string,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id":                g.ID,
		"relationship_type": "accesses-to",
		"toId":              markingDefinitionID,
	}

	queryData, err := client.Query(
		ctx,
		groupUnassignMarkingDefinitionQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[g.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: g.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

//go:embed edit_queries/group_edit_max_confidence_level.graphql
var groupEditMaxConfidenceLevelQueryString string

func (g Group) AssignMaxConfidenceLevel(
	ctx context.Context,
	client api.Client,
	confidenceLevel int,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id": g.ID,
		"input": api.EditInput{
			Key:        "group_confidence_level",
			ObjectPath: "group_confidence_level/max_confidence",
			Value:      confidenceLevel,
			Operation:  api.EditOperationReplace,
		},
	}

	queryData, err := client.Query(
		ctx,
		groupEditMaxConfidenceLevelQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[g.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: g.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

//go:embed edit_queries/group_assign_role.graphql
var groupAssignRoleQueryString string

func (g Group) AssignRole(
	ctx context.Context,
	client api.Client,
	roleID string,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id": g.ID,
		"input": entity.StixCoreRelationshipAddInput{
			RelationshipType: "has-role",
			ToID:             roleID,
		},
	}

	queryData, err := client.Query(
		ctx,
		groupAssignRoleQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[g.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: g.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

//go:embed edit_queries/group_unassign_role.graphql
var groupUnassignRoleQueryString string

func (g Group) UnassignRole(
	ctx context.Context,
	client api.Client,
	roleID string,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id":                g.ID,
		"relationship_type": "has-role",
		"toId":              roleID,
	}

	queryData, err := client.Query(
		ctx,
		groupUnassignRoleQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[g.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: g.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

//go:embed edit_queries/group_set_option.graphql
var groupSetOptionQueryString string

type GroupOption string

const (
	GroupOptionAutoNewMarking     GroupOption = "auto_new_marking"
	GroupOptionDefaultAssignation GroupOption = "default_assignation"
)

func (g Group) SetOption(
	ctx context.Context,
	client api.Client,
	option GroupOption,
	value bool,
) (map[string]any, error) {
	inputVars := map[string]any{
		"id": g.ID,
		"input": api.EditInput{
			Key:       string(option),
			Value:     value,
			Operation: api.EditOperationReplace,
		},
	}

	queryData, err := client.Query(
		ctx,
		groupSetOptionQueryString,
		inputVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to edit entity: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[g.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: g.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

// Capability utils

type CapabilityNotFoundError struct {
	name string
}

func (err CapabilityNotFoundError) Error() string { return "capability not found: " + err.name }

func (c Capability) IDsByNames(ctx context.Context, client api.Client, names []Capabilities) ([]string, error) {
	caps, err := api.List[Capability](
		ctx,
		client,
		"id, name",
		true,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot list capabilities: %w", err)
	}

	capabilities := []struct {
		ID   string `mapstructure:"id"`
		Name string `mapstructure:"name"`
	}{}
	if err := mapstructure.Decode(caps, &capabilities); err != nil {
		return nil, fmt.Errorf("unable to decode returned capabilities: %w", err)
	}

	ids := make([]string, len(names))

	for i, name := range names {
		found := false
	inner:
		for _, cap := range capabilities {
			if string(name) == cap.Name {
				found = true
				ids[i] = cap.ID

				break inner
			}
		}

		if !found {
			return nil, CapabilityNotFoundError{name: string(name)}
		}
	}

	return ids, nil
}

type Capabilities string

const (
	CapabilitiesConnectorAPI                         Capabilities = "CONNECTORAPI"
	CapabilitiesKnowledge                            Capabilities = "KNOWLEDGE"
	CapabilitiesKnowledgeKNUpdate                    Capabilities = "KNOWLEDGE_KNUPDATE"
	CapabilitiesKnowledgeKNParticipate               Capabilities = "KNOWLEDGE_KNPARTICIPATE"
	CapabilitiesKnowledgeKNUpdateKNDelete            Capabilities = "KNOWLEDGE_KNUPDATE_KNDELETE"
	CapabilitiesKnowledgeKNUpdateKNOrgAreStrict      Capabilities = "KNOWLEDGE_KNUPDATE_KNORGARESTRICT"
	CapabilitiesKnowledgeKNUpdateKNManageAuthMembers Capabilities = "KNOWLEDGE_KNUPDATE_KNMANAGEAUTHMEMBERS"
	CapabilitiesKnowledgeKNUpload                    Capabilities = "KNOWLEDGE_KNUPLOAD"
	CapabilitiesKnowledgeKNAskImport                 Capabilities = "KNOWLEDGE_KNASKIMPORT"
	CapabilitiesKnowledgeKNGetExport                 Capabilities = "KNOWLEDGE_KNGETEXPORT"
	CapabilitiesKnowledgeKNGetExportKNAskExport      Capabilities = "KNOWLEDGE_KNGETEXPORT_KNASKEXPORT"
	CapabilitiesKnowledgeKNEnrichment                Capabilities = "KNOWLEDGE_KNENRICHMENT"
	CapabilitiesExplore                              Capabilities = "EXPLORE"
	CapabilitiesExploreEXUpdate                      Capabilities = "EXPLORE_EXUPDATE"
	CapabilitiesExploreEXUpdateEXDelete              Capabilities = "EXPLORE_EXUPDATE_EXDELETE"
	CapabilitiesExploreEXUpdatePublish               Capabilities = "EXPLORE_EXUPDATE_PUBLISH"
	CapabilitiesModules                              Capabilities = "MODULES"
	CapabilitiesModulesModManage                     Capabilities = "MODULES_MODMANAGE"
	CapabilitiesSettings                             Capabilities = "SETTINGS"
	CapabilitiesSettingsSetAccesses                  Capabilities = "SETTINGS_SETACCESSES"
	CapabilitiesSettingsSetMarkings                  Capabilities = "SETTINGS_SETMARKINGS"
	CapabilitiesSettingsSetLables                    Capabilities = "SETTINGS_SETLABELS"
	CapabilitiesTAXIIAPISetCollections               Capabilities = "TAXIIAPI_SETCOLLECTIONS"
	CapabilitiesTAXIIAPISetCSVMappers                Capabilities = "TAXIIAPI_SETCSVMAPPERS"
	CapabilitiesVirtualOrganizationAdmin             Capabilities = "VIRTUAL_ORGANIZATION_ADMIN"
)

// SubType utils

func (s SubType) EditResponseField() string { return "subTypeEdit" }

//go:embed edit_queries/subtype_set_status_in_workflow.graphql
var subTypeSetStatusInWorkflowQueryString string

func (s SubType) SetStatusInWorkFlow(
	ctx context.Context,
	client api.Client,
	workflowType, statusTemplateID string,
	statusTemplateOrder int,
) (map[string]any, error) {
	queryData, err := client.Query(
		ctx,
		subTypeSetStatusInWorkflowQueryString,
		map[string]any{
			"id": workflowType,
			"input": map[string]any{
				"template_id": statusTemplateID,
				"order":       statusTemplateOrder,
			},
		})
	if err != nil {
		return nil, fmt.Errorf("cannot edit SubType: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[s.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: s.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

//go:embed edit_queries/subtype_unset_status_in_workflow.graphql
var subTypeUnsetStatusInWorkflowQueryString string

func (s SubType) UnsetStatusInWorkFlow(
	ctx context.Context,
	client api.Client,
	workflowType, statusTemplateID string,
) (map[string]any, error) {
	queryData, err := client.Query(
		ctx,
		subTypeUnsetStatusInWorkflowQueryString,
		map[string]any{
			"id": workflowType,
			"input": map[string]any{
				"status_id": statusTemplateID,
			},
		})
	if err != nil {
		return nil, fmt.Errorf("cannot edit SubType: %w", err)
	}

	// Processing the response: Expected struct: {"query_name":map[string]any}
	resp, ok := queryData[s.EditResponseField()]
	if !ok {
		return nil, api.MissingFieldError{FieldName: s.EditResponseField()}
	}

	finalMap := map[string]any{}
	if err := mapstructure.Decode(resp, &finalMap); err != nil {
		return nil, fmt.Errorf("failed to retrieve entity map: %w", err)
	}

	return finalMap, nil
}

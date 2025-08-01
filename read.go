// Code generated by '/tools/gocti_type_generator' for OpenCTI version 6.7.8 - DO NOT EDIT.

package gocti

import (
	"context"

	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/entity"
	"github.com/weisshorn-cyd/gocti/system"
)

// Helper methods declarations

func (client *OpenCTIAPIClient) ReadAttackPattern(
	ctx context.Context,
	customAttributes, id string,
) (entity.AttackPattern, error) {
	return api.StructuredRead[entity.AttackPattern, entity.AttackPattern](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadCampaign(
	ctx context.Context,
	customAttributes, id string,
) (entity.Campaign, error) {
	return api.StructuredRead[entity.Campaign, entity.Campaign](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadCaseIncident(
	ctx context.Context,
	customAttributes, id string,
) (entity.CaseIncident, error) {
	return api.StructuredRead[entity.CaseIncident, entity.CaseIncident](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadCaseRfi(
	ctx context.Context,
	customAttributes, id string,
) (entity.CaseRfi, error) {
	return api.StructuredRead[entity.CaseRfi, entity.CaseRfi](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadCaseRft(
	ctx context.Context,
	customAttributes, id string,
) (entity.CaseRft, error) {
	return api.StructuredRead[entity.CaseRft, entity.CaseRft](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadCaseTemplate(
	ctx context.Context,
	customAttributes, id string,
) (system.CaseTemplate, error) {
	return api.StructuredRead[system.CaseTemplate, system.CaseTemplate](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadChannel(
	ctx context.Context,
	customAttributes, id string,
) (entity.Channel, error) {
	return api.StructuredRead[entity.Channel, entity.Channel](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadCourseOfAction(
	ctx context.Context,
	customAttributes, id string,
) (entity.CourseOfAction, error) {
	return api.StructuredRead[entity.CourseOfAction, entity.CourseOfAction](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadDataComponent(
	ctx context.Context,
	customAttributes, id string,
) (entity.DataComponent, error) {
	return api.StructuredRead[entity.DataComponent, entity.DataComponent](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadDataSource(
	ctx context.Context,
	customAttributes, id string,
) (entity.DataSource, error) {
	return api.StructuredRead[entity.DataSource, entity.DataSource](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadEvent(
	ctx context.Context,
	customAttributes, id string,
) (entity.Event, error) {
	return api.StructuredRead[entity.Event, entity.Event](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadExternalReference(
	ctx context.Context,
	customAttributes, id string,
) (entity.ExternalReference, error) {
	return api.StructuredRead[entity.ExternalReference, entity.ExternalReference](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadFeedback(
	ctx context.Context,
	customAttributes, id string,
) (entity.Feedback, error) {
	return api.StructuredRead[entity.Feedback, entity.Feedback](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadGroup(
	ctx context.Context,
	customAttributes, id string,
) (system.Group, error) {
	return api.StructuredRead[system.Group, system.Group](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadGrouping(
	ctx context.Context,
	customAttributes, id string,
) (entity.Grouping, error) {
	return api.StructuredRead[entity.Grouping, entity.Grouping](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadIdentity(
	ctx context.Context,
	customAttributes, id string,
) (entity.Identity, error) {
	return api.StructuredRead[entity.Identity, entity.Identity](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadIncident(
	ctx context.Context,
	customAttributes, id string,
) (entity.Incident, error) {
	return api.StructuredRead[entity.Incident, entity.Incident](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadIndicator(
	ctx context.Context,
	customAttributes, id string,
) (entity.Indicator, error) {
	return api.StructuredRead[entity.Indicator, entity.Indicator](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadInfrastructure(
	ctx context.Context,
	customAttributes, id string,
) (entity.Infrastructure, error) {
	return api.StructuredRead[entity.Infrastructure, entity.Infrastructure](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadIntrusionSet(
	ctx context.Context,
	customAttributes, id string,
) (entity.IntrusionSet, error) {
	return api.StructuredRead[entity.IntrusionSet, entity.IntrusionSet](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadKillChainPhase(
	ctx context.Context,
	customAttributes, id string,
) (entity.KillChainPhase, error) {
	return api.StructuredRead[entity.KillChainPhase, entity.KillChainPhase](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadLabel(
	ctx context.Context,
	customAttributes, id string,
) (entity.Label, error) {
	return api.StructuredRead[entity.Label, entity.Label](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadLanguage(
	ctx context.Context,
	customAttributes, id string,
) (entity.Language, error) {
	return api.StructuredRead[entity.Language, entity.Language](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadLocation(
	ctx context.Context,
	customAttributes, id string,
) (entity.Location, error) {
	return api.StructuredRead[entity.Location, entity.Location](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadMalware(
	ctx context.Context,
	customAttributes, id string,
) (entity.Malware, error) {
	return api.StructuredRead[entity.Malware, entity.Malware](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadMalwareAnalysis(
	ctx context.Context,
	customAttributes, id string,
) (entity.MalwareAnalysis, error) {
	return api.StructuredRead[entity.MalwareAnalysis, entity.MalwareAnalysis](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadMarkingDefinition(
	ctx context.Context,
	customAttributes, id string,
) (entity.MarkingDefinition, error) {
	return api.StructuredRead[entity.MarkingDefinition, entity.MarkingDefinition](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadNarrative(
	ctx context.Context,
	customAttributes, id string,
) (entity.Narrative, error) {
	return api.StructuredRead[entity.Narrative, entity.Narrative](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadNote(
	ctx context.Context,
	customAttributes, id string,
) (entity.Note, error) {
	return api.StructuredRead[entity.Note, entity.Note](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadObservedData(
	ctx context.Context,
	customAttributes, id string,
) (entity.ObservedData, error) {
	return api.StructuredRead[entity.ObservedData, entity.ObservedData](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadOpinion(
	ctx context.Context,
	customAttributes, id string,
) (entity.Opinion, error) {
	return api.StructuredRead[entity.Opinion, entity.Opinion](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadReport(
	ctx context.Context,
	customAttributes, id string,
) (entity.Report, error) {
	return api.StructuredRead[entity.Report, entity.Report](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadRole(
	ctx context.Context,
	customAttributes, id string,
) (system.Role, error) {
	return api.StructuredRead[system.Role, system.Role](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadStatusTemplate(
	ctx context.Context,
	customAttributes, id string,
) (system.StatusTemplate, error) {
	return api.StructuredRead[system.StatusTemplate, system.StatusTemplate](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadStixCoreObject(
	ctx context.Context,
	customAttributes, id string,
) (entity.StixCoreObject, error) {
	return api.StructuredRead[entity.StixCoreObject, entity.StixCoreObject](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadStixCoreRelationship(
	ctx context.Context,
	customAttributes, id string,
) (entity.StixCoreRelationship, error) {
	return api.StructuredRead[entity.StixCoreRelationship, entity.StixCoreRelationship](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadStixCyberObservable(
	ctx context.Context,
	customAttributes, id string,
) (entity.StixCyberObservable, error) {
	return api.StructuredRead[entity.StixCyberObservable, entity.StixCyberObservable](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadStixDomainObject(
	ctx context.Context,
	customAttributes, id string,
) (entity.StixDomainObject, error) {
	return api.StructuredRead[entity.StixDomainObject, entity.StixDomainObject](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadSubType(
	ctx context.Context,
	customAttributes, id string,
) (system.SubType, error) {
	return api.StructuredRead[system.SubType, system.SubType](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadTask(
	ctx context.Context,
	customAttributes, id string,
) (entity.Task, error) {
	return api.StructuredRead[entity.Task, entity.Task](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadTaskTemplate(
	ctx context.Context,
	customAttributes, id string,
) (system.TaskTemplate, error) {
	return api.StructuredRead[system.TaskTemplate, system.TaskTemplate](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadThreatActor(
	ctx context.Context,
	customAttributes, id string,
) (entity.ThreatActor, error) {
	return api.StructuredRead[entity.ThreatActor, entity.ThreatActor](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadThreatActorGroup(
	ctx context.Context,
	customAttributes, id string,
) (entity.ThreatActorGroup, error) {
	return api.StructuredRead[entity.ThreatActorGroup, entity.ThreatActorGroup](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadThreatActorIndividual(
	ctx context.Context,
	customAttributes, id string,
) (entity.ThreatActorIndividual, error) {
	return api.StructuredRead[entity.ThreatActorIndividual, entity.ThreatActorIndividual](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadTool(
	ctx context.Context,
	customAttributes, id string,
) (entity.Tool, error) {
	return api.StructuredRead[entity.Tool, entity.Tool](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadUser(
	ctx context.Context,
	customAttributes, id string,
) (system.User, error) {
	return api.StructuredRead[system.User, system.User](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadVocabulary(
	ctx context.Context,
	customAttributes, id string,
) (entity.Vocabulary, error) {
	return api.StructuredRead[entity.Vocabulary, entity.Vocabulary](
		ctx, client, customAttributes, id,
	)
}

func (client *OpenCTIAPIClient) ReadVulnerability(
	ctx context.Context,
	customAttributes, id string,
) (entity.Vulnerability, error) {
	return api.StructuredRead[entity.Vulnerability, entity.Vulnerability](
		ctx, client, customAttributes, id,
	)
}

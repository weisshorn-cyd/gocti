// Code generated by '/tools/gocti_type_generator' for OpenCTI version 6.7.8 - DO NOT EDIT.

package gocti

import (
	"context"

	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/entity"
	"github.com/weisshorn-cyd/gocti/graphql"
	"github.com/weisshorn-cyd/gocti/list"
	"github.com/weisshorn-cyd/gocti/system"
)

// Helper methods declarations

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListAttackPatterns(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.AttackPattern, error) {
	return api.StructuredList[entity.AttackPattern, entity.AttackPattern](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListCampaigns(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Campaign, error) {
	return api.StructuredList[entity.Campaign, entity.Campaign](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithFirst]
func (client *OpenCTIAPIClient) ListCapabilities(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]system.Capability, error) {
	return api.StructuredList[system.Capability, system.Capability](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListCaseIncidents(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.CaseIncident, error) {
	return api.StructuredList[entity.CaseIncident, entity.CaseIncident](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListCaseRfis(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.CaseRfi, error) {
	return api.StructuredList[entity.CaseRfi, entity.CaseRfi](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListCaseRfts(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.CaseRft, error) {
	return api.StructuredList[entity.CaseRft, entity.CaseRft](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListCaseTemplates(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]system.CaseTemplate, error) {
	return api.StructuredList[system.CaseTemplate, system.CaseTemplate](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListChannels(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Channel, error) {
	return api.StructuredList[entity.Channel, entity.Channel](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListCoursesOfAction(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.CourseOfAction, error) {
	return api.StructuredList[entity.CourseOfAction, entity.CourseOfAction](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListDataComponents(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.DataComponent, error) {
	return api.StructuredList[entity.DataComponent, entity.DataComponent](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListDataSources(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.DataSource, error) {
	return api.StructuredList[entity.DataSource, entity.DataSource](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListEvents(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Event, error) {
	return api.StructuredList[entity.Event, entity.Event](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListExternalReferences(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.ExternalReference, error) {
	return api.StructuredList[entity.ExternalReference, entity.ExternalReference](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListFeedbacks(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Feedback, error) {
	return api.StructuredList[entity.Feedback, entity.Feedback](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListGroups(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]system.Group, error) {
	return api.StructuredList[system.Group, system.Group](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListGroupings(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Grouping, error) {
	return api.StructuredList[entity.Grouping, entity.Grouping](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
//   - [list.WithTypes]
func (client *OpenCTIAPIClient) ListIdentities(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Identity, error) {
	return api.StructuredList[entity.Identity, entity.Identity](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListIncidents(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Incident, error) {
	return api.StructuredList[entity.Incident, entity.Incident](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListIndicators(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Indicator, error) {
	return api.StructuredList[entity.Indicator, entity.Indicator](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListInfrastructures(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Infrastructure, error) {
	return api.StructuredList[entity.Infrastructure, entity.Infrastructure](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListIntrusionSets(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.IntrusionSet, error) {
	return api.StructuredList[entity.IntrusionSet, entity.IntrusionSet](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListKillChainPhases(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.KillChainPhase, error) {
	return api.StructuredList[entity.KillChainPhase, entity.KillChainPhase](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListLabels(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Label, error) {
	return api.StructuredList[entity.Label, entity.Label](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListLanguages(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Language, error) {
	return api.StructuredList[entity.Language, entity.Language](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
//   - [list.WithTypes]
func (client *OpenCTIAPIClient) ListLocations(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Location, error) {
	return api.StructuredList[entity.Location, entity.Location](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListMalwares(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Malware, error) {
	return api.StructuredList[entity.Malware, entity.Malware](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListMalwareAnalyses(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.MalwareAnalysis, error) {
	return api.StructuredList[entity.MalwareAnalysis, entity.MalwareAnalysis](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListMarkingDefinitions(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.MarkingDefinition, error) {
	return api.StructuredList[entity.MarkingDefinition, entity.MarkingDefinition](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListNarratives(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Narrative, error) {
	return api.StructuredList[entity.Narrative, entity.Narrative](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListNotes(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Note, error) {
	return api.StructuredList[entity.Note, entity.Note](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListObservedDatas(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.ObservedData, error) {
	return api.StructuredList[entity.ObservedData, entity.ObservedData](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListOpinions(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Opinion, error) {
	return api.StructuredList[entity.Opinion, entity.Opinion](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListReports(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Report, error) {
	return api.StructuredList[entity.Report, entity.Report](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListRoles(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]system.Role, error) {
	return api.StructuredList[system.Role, system.Role](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListStatusTemplates(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]system.StatusTemplate, error) {
	return api.StructuredList[system.StatusTemplate, system.StatusTemplate](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithTypes]
func (client *OpenCTIAPIClient) ListStixCoreObjects(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.StixCoreObject, error) {
	return api.StructuredList[entity.StixCoreObject, entity.StixCoreObject](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithConfidences]
//   - [list.WithDynamicFrom]
//   - [list.WithDynamicTo]
//   - [list.WithElementWithTargetTypes]
//   - [list.WithEndDate]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithFirstSeenStart]
//   - [list.WithFirstSeenStop]
//   - [list.WithFromIDList]
//   - [list.WithFromOrToIDList]
//   - [list.WithFromRole]
//   - [list.WithFromTypes]
//   - [list.WithLastSeenStart]
//   - [list.WithLastSeenStop]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithRelationshipType]
//   - [list.WithSearch]
//   - [list.WithStartDate]
//   - [list.WithStartTimeStart]
//   - [list.WithStartTimeStop]
//   - [list.WithStix]
//   - [list.WithStopTimeStart]
//   - [list.WithStopTimeStop]
//   - [list.WithToIDList]
//   - [list.WithToRole]
//   - [list.WithToTypes]
func (client *OpenCTIAPIClient) ListStixCoreRelationships(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.StixCoreRelationship, error) {
	return api.StructuredList[entity.StixCoreRelationship, entity.StixCoreRelationship](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
//   - [list.WithTypes]
func (client *OpenCTIAPIClient) ListStixCyberObservables(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.StixCyberObservable, error) {
	return api.StructuredList[entity.StixCyberObservable, entity.StixCyberObservable](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithTypes]
func (client *OpenCTIAPIClient) ListStixDomainObjects(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.StixDomainObject, error) {
	return api.StructuredList[entity.StixDomainObject, entity.StixDomainObject](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFirst]
//   - [list.WithIncludeParents]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithType]
func (client *OpenCTIAPIClient) ListSubTypes(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]system.SubType, error) {
	return api.StructuredList[system.SubType, system.SubType](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListTasks(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Task, error) {
	return api.StructuredList[entity.Task, entity.Task](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListTaskTemplates(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]system.TaskTemplate, error) {
	return api.StructuredList[system.TaskTemplate, system.TaskTemplate](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListThreatActors(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.ThreatActor, error) {
	return api.StructuredList[entity.ThreatActor, entity.ThreatActor](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListThreatActorGroups(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.ThreatActorGroup, error) {
	return api.StructuredList[entity.ThreatActorGroup, entity.ThreatActorGroup](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListThreatActorIndividuals(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.ThreatActorIndividual, error) {
	return api.StructuredList[entity.ThreatActorIndividual, entity.ThreatActorIndividual](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListTools(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Tool, error) {
	return api.StructuredList[entity.Tool, entity.Tool](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListUsers(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]system.User, error) {
	return api.StructuredList[system.User, system.User](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithCategory]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
func (client *OpenCTIAPIClient) ListVocabularies(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Vocabulary, error) {
	return api.StructuredList[entity.Vocabulary, entity.Vocabulary](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

// Available [list.Option] parameters:
//   - [list.WithAfter]
//   - [list.WithFilters]
//   - [list.WithFirst]
//   - [list.WithOrderBy]
//   - [list.WithOrderMode]
//   - [list.WithSearch]
//   - [list.WithToStix]
func (client *OpenCTIAPIClient) ListVulnerabilities(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *graphql.PageInfo,
	opts ...list.Option,
) ([]entity.Vulnerability, error) {
	return api.StructuredList[entity.Vulnerability, entity.Vulnerability](
		ctx, client, customAttributes, getAll, pageInfo, opts...,
	)
}

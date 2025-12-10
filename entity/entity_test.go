package entity_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	"github.com/weisshorn-cyd/gocti"
	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/entity"
	"github.com/weisshorn-cyd/gocti/graphql"
	"github.com/weisshorn-cyd/gocti/list"
)

const (
	timeout = 3 * time.Second
)

type config struct {
	URL   string `envconfig:"URL"   required:"true"`
	Token string `envconfig:"TOKEN" required:"true"`
}

func loadConfig(t *testing.T) config {
	t.Helper()

	config := config{}
	if err := envconfig.Process("OPENCTI", &config); err != nil {
		t.Fatalf("unable to load test config: %v", err)
	}

	return config
}

// args is a struct that holds all functions and data to run
// a complete test of an entitie's related queries.
type args struct {
	createFunc func(
		ctx context.Context,
		client api.Client,
		customAttributes string,
		input api.AddInput,
	) (map[string]any, error)
	readFunc func(
		ctx context.Context,
		client api.Client,
		customAttributes, id string,
	) (map[string]any, error)
	listFunc func(
		ctx context.Context,
		client api.Client,
		customAttributes string,
		getAll bool,
		pageInfo *graphql.PageInfo,
		opts ...list.Option,
	) ([]map[string]any, error)
	deleteFunc func(
		ctx context.Context,
		id string,
	) (string, error)
	input   api.AddInput
	orderBy string
}

//nolint:maintidx // Triggers on the high number of cases, but they are only small tests.
func TestEntity(t *testing.T) {
	cfg := loadConfig(t)

	client, err := gocti.NewOpenCTIAPIClient(
		cfg.URL, cfg.Token,
		gocti.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	)
	if err != nil {
		t.Fatalf("cannot create client: %v", err)
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test AttackPattern",
			args: args{
				createFunc: api.Create[entity.AttackPattern],
				readFunc:   api.Read[entity.AttackPattern],
				listFunc:   api.List[entity.AttackPattern],
				deleteFunc: client.DeleteAttackPattern,
				input: entity.AttackPatternAddInput{
					Name: "Test AttackPattern",
				},
			},
		}, {
			name: "Test Campaign",
			args: args{
				createFunc: api.Create[entity.Campaign],
				readFunc:   api.Read[entity.Campaign],
				listFunc:   api.List[entity.Campaign],
				deleteFunc: client.DeleteCampaign,
				input: entity.CampaignAddInput{
					Name: "Test Campaign",
				},
			},
		}, {
			name: "Test CaseIncident",
			args: args{
				createFunc: api.Create[entity.CaseIncident],
				readFunc:   api.Read[entity.CaseIncident],
				listFunc:   api.List[entity.CaseIncident],
				deleteFunc: client.DeleteCaseIncident,
				input: entity.CaseIncidentAddInput{
					Name: "Test CaseIncident",
				},
			},
		}, {
			name: "Test CaseRfi",
			args: args{
				createFunc: api.Create[entity.CaseRfi],
				readFunc:   api.Read[entity.CaseRfi],
				listFunc:   api.List[entity.CaseRfi],
				deleteFunc: client.DeleteCaseRfi,
				input: entity.CaseRfiAddInput{
					Name: "Test CaseRfi",
				},
			},
		}, {
			name: "Test CaseRft",
			args: args{
				createFunc: api.Create[entity.CaseRft],
				readFunc:   api.Read[entity.CaseRft],
				listFunc:   api.List[entity.CaseRft],
				deleteFunc: client.DeleteCaseRft,
				input: entity.CaseRftAddInput{
					Name: "Test CaseRft",
				},
			},
		}, {
			name: "Test Channel",
			args: args{
				createFunc: api.Create[entity.Channel],
				readFunc:   api.Read[entity.Channel],
				listFunc:   api.List[entity.Channel],
				deleteFunc: client.DeleteChannel,
				input: entity.ChannelAddInput{
					Name: "Test Channel",
				},
			},
		}, {
			name: "Test CourseOfAction",
			args: args{
				createFunc: api.Create[entity.CourseOfAction],
				readFunc:   api.Read[entity.CourseOfAction],
				listFunc:   api.List[entity.CourseOfAction],
				deleteFunc: client.DeleteCourseOfAction,
				input: entity.CourseOfActionAddInput{
					Name: "Test CourseOfAction",
				},
			},
		}, {
			name: "Test DataComponent",
			args: args{
				createFunc: api.Create[entity.DataComponent],
				readFunc:   api.Read[entity.DataComponent],
				listFunc:   api.List[entity.DataComponent],
				deleteFunc: client.DeleteDataComponent,
				input: entity.DataComponentAddInput{
					Name: "Test DataComponent",
				},
			},
		}, {
			name: "Test DataSource",
			args: args{
				createFunc: api.Create[entity.DataSource],
				readFunc:   api.Read[entity.DataSource],
				listFunc:   api.List[entity.DataSource],
				deleteFunc: client.DeleteDataSource,
				input: entity.DataSourceAddInput{
					Name: "Test DataSource",
				},
			},
		}, {
			name: "Test Event",
			args: args{
				createFunc: api.Create[entity.Event],
				readFunc:   api.Read[entity.Event],
				listFunc:   api.List[entity.Event],
				deleteFunc: client.DeleteEvent,
				input: entity.EventAddInput{
					Name: "Test Event",
				},
			},
		}, {
			name: "Test ExternalReference",
			args: args{
				createFunc: api.Create[entity.ExternalReference],
				readFunc:   api.Read[entity.ExternalReference],
				listFunc:   api.List[entity.ExternalReference],
				deleteFunc: client.DeleteExternalReference,
				input: entity.ExternalReferenceAddInput{
					SourceName: "source_name",
					URL:        "url",
				},
			},
		}, {
			name: "Test Feedback",
			args: args{
				createFunc: api.Create[entity.Feedback],
				readFunc:   api.Read[entity.Feedback],
				listFunc:   api.List[entity.Feedback],
				deleteFunc: client.DeleteFeedback,
				input: entity.FeedbackAddInput{
					Name: "Test Feedback",
				},
			},
		}, {
			name: "Test Grouping",
			args: args{
				createFunc: api.Create[entity.Grouping],
				readFunc:   api.Read[entity.Grouping],
				listFunc:   api.List[entity.Grouping],
				deleteFunc: client.DeleteGrouping,
				input: entity.GroupingAddInput{
					Name:    "Test Grouping",
					Context: graphql.GroupingContextMalwareAnalysis,
				},
			},
		}, {
			name: "Test Identity",
			args: args{
				createFunc: api.Create[entity.Identity],
				readFunc:   api.Read[entity.Identity],
				listFunc:   api.List[entity.Identity],
				deleteFunc: client.DeleteIdentity,
				input: entity.IdentityAddInput{
					Type: graphql.IdentityTypeSector,
					Name: "Test Identity",
				},
			},
		}, {
			name: "Test Incident",
			args: args{
				createFunc: api.Create[entity.Incident],
				readFunc:   api.Read[entity.Incident],
				listFunc:   api.List[entity.Incident],
				deleteFunc: client.DeleteIncident,
				input: entity.IncidentAddInput{
					Name: "Test Incident",
				},
			},
		}, {
			name: "Test Indicator",
			args: args{
				createFunc: api.Create[entity.Indicator],
				readFunc:   api.Read[entity.Indicator],
				listFunc:   api.List[entity.Indicator],
				deleteFunc: client.DeleteIndicator,
				input: entity.IndicatorAddInput{
					PatternType: graphql.PatternTypeStix,
					Pattern:     "[domain-name:value = 'www.test.test']",
					Name:        "Test Indicator",
				},
			},
		}, {
			name: "Test Infrastructure",
			args: args{
				createFunc: api.Create[entity.Infrastructure],
				readFunc:   api.Read[entity.Infrastructure],
				listFunc:   api.List[entity.Infrastructure],
				deleteFunc: client.DeleteInfrastructure,
				input: entity.InfrastructureAddInput{
					Name: "Test Infrastructure",
				},
			},
		}, {
			name: "Test IntrusionSet",
			args: args{
				createFunc: api.Create[entity.IntrusionSet],
				readFunc:   api.Read[entity.IntrusionSet],
				listFunc:   api.List[entity.IntrusionSet],
				deleteFunc: client.DeleteIntrusionSet,
				input: entity.IntrusionSetAddInput{
					Name: "Test IntrusionSet",
				},
			},
		}, {
			name: "Test KillChainPhase",
			args: args{
				createFunc: api.Create[entity.KillChainPhase],
				readFunc:   api.Read[entity.KillChainPhase],
				listFunc:   api.List[entity.KillChainPhase],
				deleteFunc: client.DeleteKillChainPhase,
				input: entity.KillChainPhaseAddInput{
					KillChainName: "Test Kill Chain",
					PhaseName:     "Test Phase",
					XOpenctiOrder: 1,
				},
			},
		}, {
			name: "Test Label",
			args: args{
				createFunc: api.Create[entity.Label],
				readFunc:   api.Read[entity.Label],
				listFunc:   api.List[entity.Label],
				deleteFunc: client.DeleteLabel,
				input: entity.LabelAddInput{
					Value: "Test Label",
				},
			},
		}, {
			name: "Test Language",
			args: args{
				createFunc: api.Create[entity.Language],
				readFunc:   api.Read[entity.Language],
				listFunc:   api.List[entity.Language],
				deleteFunc: client.DeleteLanguage,
				input: entity.LanguageAddInput{
					Name: "Test Language",
				},
			},
		}, {
			name: "Test Location",
			args: args{
				createFunc: api.Create[entity.Location],
				readFunc:   api.Read[entity.Location],
				listFunc:   api.List[entity.Location],
				deleteFunc: client.DeleteLocation,
				input: entity.LocationAddInput{
					Type: "Region",
					Name: "Test Location",
				},
			},
		}, {
			name: "Test Malware",
			args: args{
				createFunc: api.Create[entity.Malware],
				readFunc:   api.Read[entity.Malware],
				listFunc:   api.List[entity.Malware],
				deleteFunc: client.DeleteMalware,
				input: entity.MalwareAddInput{
					Name: "Test Malware",
				},
			},
		}, {
			name: "Test MalwareAnalysis",
			args: args{
				createFunc: api.Create[entity.MalwareAnalysis],
				readFunc:   api.Read[entity.MalwareAnalysis],
				listFunc:   api.List[entity.MalwareAnalysis],
				deleteFunc: client.DeleteMalwareAnalysis,
				input: entity.MalwareAnalysisAddInput{
					Product:    "Test Product",
					ResultName: "Result",
				},
				orderBy: "creator",
			},
		}, {
			name: "Test MarkingDefinition",
			args: args{
				createFunc: api.Create[entity.MarkingDefinition],
				readFunc:   api.Read[entity.MarkingDefinition],
				listFunc:   api.List[entity.MarkingDefinition],
				deleteFunc: client.DeleteMarkingDefinition,
				input: entity.MarkingDefinitionAddInput{
					DefinitionType: "TEST",
					Definition:     "TEST",
					XOpenctiOrder:  2,
				},
			},
		}, {
			name: "Test Narrative",
			args: args{
				createFunc: api.Create[entity.Narrative],
				readFunc:   api.Read[entity.Narrative],
				listFunc:   api.List[entity.Narrative],
				deleteFunc: client.DeleteNarrative,
				input: entity.NarrativeAddInput{
					Name: "Test Narrative",
				},
			},
		}, {
			name: "Test Note",
			args: args{
				createFunc: api.Create[entity.Note],
				readFunc:   api.Read[entity.Note],
				listFunc:   api.List[entity.Note],
				deleteFunc: client.DeleteNote,
				input: entity.NoteAddInput{
					Content: "content",
				},
			},
		}, {
			name: "Test Opinion",
			args: args{
				createFunc: api.Create[entity.Opinion],
				readFunc:   api.Read[entity.Opinion],
				listFunc:   api.List[entity.Opinion],
				deleteFunc: client.DeleteOpinion,
				input: entity.OpinionAddInput{
					Opinion: graphql.OpinionTypeNeutral,
				},
			},
		}, {
			name: "Test Report",
			args: args{
				createFunc: api.Create[entity.Report],
				readFunc:   api.Read[entity.Report],
				listFunc:   api.List[entity.Report],
				deleteFunc: client.DeleteReport,
				input: entity.ReportAddInput{
					Name:      "Test Report",
					Published: &time.Time{},
				},
			},
		}, {
			name: "Test StixCoreObject",
			args: args{
				createFunc: nil,
				readFunc:   nil,
				listFunc:   api.List[entity.StixCoreObject],
				deleteFunc: nil,
				input:      nil,
			},
		}, {
			name: "Test StixCyberObservable",
			args: args{
				createFunc: api.Create[entity.StixCyberObservable],
				readFunc:   api.Read[entity.StixCyberObservable],
				listFunc:   api.List[entity.StixCyberObservable],
				deleteFunc: client.DeleteStixCyberObservable,
				input: entity.StixCyberObservableAddInput{
					Type: graphql.StixCyberObservableTypeIPV4Addr,
					IPv4Addr: graphql.IPv4AddrAddInput{
						Value: "192.0.2.43",
					},
				},
			},
		}, {
			name: "Test StixDomainObject",
			args: args{
				createFunc: api.Create[entity.StixDomainObject],
				readFunc:   api.Read[entity.StixDomainObject],
				listFunc:   api.List[entity.StixDomainObject],
				deleteFunc: client.DeleteStixDomainObject,
				input: entity.StixDomainObjectAddInput{
					Name: "Test StixDomainObject",
					Type: "Report",
				},
			},
		}, {
			name: "Test Task",
			args: args{
				createFunc: api.Create[entity.Task],
				readFunc:   api.Read[entity.Task],
				listFunc:   api.List[entity.Task],
				deleteFunc: client.DeleteTask,
				input: entity.TaskAddInput{
					Name: "Test Task",
				},
			},
		}, {
			name: "Test ThreatActor",
			args: args{
				createFunc: nil,
				readFunc:   api.Read[entity.ThreatActor],
				listFunc:   api.List[entity.ThreatActor],
				deleteFunc: nil,
				input:      nil,
			},
		}, {
			name: "Test ThreatActorGroup",
			args: args{
				createFunc: api.Create[entity.ThreatActorGroup],
				readFunc:   nil,
				listFunc:   nil,
				deleteFunc: client.DeleteThreatActorGroup,
				input: entity.ThreatActorGroupAddInput{
					Name: "Test ThreatActorGroup",
				},
			},
		}, {
			name: "Test ThreatActorIndividual",
			args: args{
				createFunc: api.Create[entity.ThreatActorIndividual],
				readFunc:   nil,
				listFunc:   nil,
				deleteFunc: client.DeleteThreatActorIndividual,
				input: entity.ThreatActorIndividualAddInput{
					Name: "Test ThreatActorIndividual",
				},
			},
		}, {
			name: "Test Tool",
			args: args{
				createFunc: api.Create[entity.Tool],
				readFunc:   api.Read[entity.Tool],
				listFunc:   api.List[entity.Tool],
				deleteFunc: client.DeleteTool,
				input: entity.ToolAddInput{
					Name: "Test Tool",
				},
			},
		}, {
			name: "Test Vocabulary",
			args: args{
				createFunc: api.Create[entity.Vocabulary],
				readFunc:   api.Read[entity.Vocabulary],
				listFunc:   api.List[entity.Vocabulary],
				deleteFunc: client.DeleteVocabulary,
				input: entity.VocabularyAddInput{
					Name:     "Test Vocabulary",
					Category: "account_type_ov",
				},
				orderBy: "name",
			},
		}, {
			name: "Test Vulnerability",
			args: args{
				createFunc: api.Create[entity.Vulnerability],
				readFunc:   api.Read[entity.Vulnerability],
				listFunc:   api.List[entity.Vulnerability],
				deleteFunc: client.DeleteVulnerability,
				input: entity.VulnerabilityAddInput{
					Name: "Test Vulnerability",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id := "undefined"

			if test.args.createFunc != nil && test.args.deleteFunc != nil {
				id = createEntityAndScheduleDelete(t, client, test.args, test.name)
			}

			if test.args.readFunc != nil && id != "undefined" {
				readEntity(t, client, test.args, id, test.name)
			}

			if test.args.listFunc != nil {
				if test.args.orderBy == "" {
					test.args.orderBy = "created_at"
				}

				listEntities(t, client, test.args, test.name)
			}
		})
	}
}

func TestRelationship(t *testing.T) {
	cfg := loadConfig(t)

	client, err := gocti.NewOpenCTIAPIClient(
		cfg.URL, cfg.Token,
		gocti.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	)
	if err != nil {
		t.Fatalf("cannot create client: %v", err)
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test StixCoreRelationship",
			args: args{
				createFunc: api.Create[entity.StixCoreRelationship],
				readFunc:   api.Read[entity.StixCoreRelationship],
				listFunc:   api.List[entity.StixCoreRelationship],
				deleteFunc: client.DeleteStixCoreRelationship,
				input: &entity.StixCoreRelationshipAddInput{
					FromID:           "",
					ToID:             "",
					RelationshipType: "attributed-to",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			threatActorName := fmt.Sprintf("Test Threat Actor Individual (%s)", test.name)

			idThreatActor := createEntityAndScheduleDelete(
				t, client, args{
					createFunc: api.Create[entity.ThreatActorIndividual],
					deleteFunc: client.DeleteThreatActorIndividual,
					input: entity.ThreatActorIndividualAddInput{
						Name: threatActorName,
					},
				},
				threatActorName,
			)

			campaignName := fmt.Sprintf("Test Campaign (%s)", test.name)

			idCampaign := createEntityAndScheduleDelete(
				t, client, args{
					createFunc: api.Create[entity.Campaign],
					deleteFunc: client.DeleteCampaign,
					input: entity.CampaignAddInput{
						Name: campaignName,
					},
				},
				campaignName,
			)

			input, ok := test.args.input.(*entity.StixCoreRelationshipAddInput)
			if ok {
				input.FromID = idCampaign
				input.ToID = idThreatActor
			}

			idRelationship := createEntityAndScheduleDelete(
				t, client, args{
					createFunc: test.args.createFunc,
					deleteFunc: test.args.deleteFunc,
					input:      test.args.input,
				},
				test.name,
			)

			readEntity(
				t, client, args{
					readFunc: test.args.readFunc,
				},
				idRelationship, test.name,
			)

			listEntities(
				t, client, args{
					listFunc: test.args.listFunc,
				},
				test.name,
			)
		})
	}
}

func TestObservedData(t *testing.T) {
	cfg := loadConfig(t)

	client, err := gocti.NewOpenCTIAPIClient(
		cfg.URL, cfg.Token,
		gocti.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	)
	if err != nil {
		t.Fatalf("cannot create client: %v", err)
	}

	idCampaign := createEntityAndScheduleDelete(
		t, client, args{
			createFunc: api.Create[entity.Campaign],
			deleteFunc: client.DeleteCampaign,
			input: entity.CampaignAddInput{
				Name: "Test Campaign (ObservedData)",
			},
		},
		"Test Campaign",
	)

	idObservedData := createEntityAndScheduleDelete(
		t, client, args{
			createFunc: api.Create[entity.ObservedData],
			deleteFunc: client.DeleteObservedData,
			input: entity.ObservedDataAddInput{
				FirstObserved:  &time.Time{},
				LastObserved:   &time.Time{},
				NumberObserved: 2,
				Objects:        []string{idCampaign},
			},
		},
		"Test Observed Data",
	)

	readEntity(
		t, client, args{
			readFunc: api.Read[entity.ObservedData],
		},
		idObservedData, "Test Observed Data",
	)

	listEntities(
		t, client, args{
			listFunc: api.List[entity.ObservedData],
		},
		"Test Observed Data",
	)
}

func createEntityAndScheduleDelete(t *testing.T, client *gocti.OpenCTIAPIClient, args args, name string) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create the Entity
	createResp, err := args.createFunc(
		ctx, client, "id", args.input,
	)
	require.NoError(t, err)

	idAny, ok := createResp["id"]
	require.Truef(t, ok, "cannot get id after %s creation", name)

	id, ok := idAny.(string)
	require.Truef(t, ok, "cannot get id string after %s creation", name)

	// Setup the Entity's destruction
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_, err := args.deleteFunc(
			ctx, id,
		)
		require.NoErrorf(t, err, "failed to remove %s", name)
	})

	return id
}

func readEntity(t *testing.T, client *gocti.OpenCTIAPIClient, args args, id, name string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Test the Read function on the Label
	_, err := args.readFunc(
		ctx, client, "", id,
	)
	if err != nil {
		t.Errorf("cannot read %s: %v", name, err)
	}
}

func listEntities(t *testing.T, client *gocti.OpenCTIAPIClient, args args, name string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Test the List function on the Label
	_, err := args.listFunc(
		ctx, client, "", false, nil,
		list.WithOrderBy(args.orderBy),
	)
	if err != nil {
		t.Errorf("cannot list %s: %v", name, err)
	}
}

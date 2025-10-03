// package main_test contains examples of GoCTI usage.
//
// This example highlights some current limitations of GoCTI and how to deal
// with them. The list of observables contained in a report are not retrievable
// by default. We need to extend the Report struct and define a filter using GraphQL.
package main_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/weisshorn-cyd/gocti"
	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/entity"
	"github.com/weisshorn-cyd/gocti/graphql"
	"github.com/weisshorn-cyd/gocti/list"
)

// We need to extend the entity.Report struct with the Objects attribute we want
// to retrieve.
type ExtendedReport struct {
	entity.Report

	Objects struct {
		Edges []struct {
			Node Observable `gocti:"node"`
		} `gocti:"edges"`
	} `gocti:"objects"`
}

type Observable struct {
	ID      string `gocti:"id"`
	Value   string `gocti:"observable_value"`
	ModDate string `gocti:"updated_at"`
}

// reportQueryAttributes are the custom attributes matching the above Report and
// Observable structures.
// This parameter field must be filled with the result of FilterGroup.ToGraphQL():
// -------------------------------------------- v.
const reportQueryAttributes = `objects(filters: %s){
		edges {
			node {
				... on StixCyberObservable {
					updated_at
					observable_value
					id
				}
			}
		}
	}`

func Example_listReportIPs() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	// Get config from env
	cfg := struct {
		URL   string `envconfig:"URL"   required:"true"`
		Token string `envconfig:"TOKEN" required:"true"`
	}{}
	if err := envconfig.Process("OPENCTI", &cfg); err != nil {
		logger.Error("loading config", "error", err)
		os.Exit(1)
	}

	// Create client
	client, err := gocti.NewOpenCTIAPIClient(
		cfg.URL,
		cfg.Token,
		gocti.WithHealthCheck(),
		gocti.WithLogger(logger),
	)
	if err != nil {
		logger.Error("creating client", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Provision
	callback := provision(ctx, client, logger)
	defer callback()

	// List any number of label names.
	labels := []any{"c2-fingerprints", "ioc"}
	// Only selecting observables later or equal to that date.
	observableCutoffDate := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.Now().Location())
	// List of observable types to filter.
	observableTypes := []any{"IPv4-Addr", "IPv6-Addr"}

	// Get selected label(s)' ID
	labelIDs, err := getLabelIDs(ctx, client, labels)
	if err != nil {
		logger.Error("retrieving label ids", "error", err)
	}

	fmt.Printf("Found %d entries with the requested label(s)\n", len(labelIDs))

	// Get reports tagged with the selected label(s)
	// The contained observables are filtered to only be of the chosen type(s) and
	// to have been modified no earlier than the cutoff date
	reports, err := getReports(ctx, client, labelIDs, observableTypes, observableCutoffDate)
	if err != nil {
		logger.Error("retrieving labeled reports", "error", err)
	}

	fmt.Printf("Found %d Report(s)\n", len(reports))

	// Extract the observables from the reports and remove duplicates
	observables := getUniqueObservables(reports)

	fmt.Printf("Found %d Observable(s) total\n", len(observables))

	// Output:
	// Found 1 entries with the requested label(s)
	// Found 1 Report(s)
	// Found 1 Observable(s) total
}

// nolint: cyclop, nolintlint // ok for a usage example
func provision(ctx context.Context, client *gocti.OpenCTIAPIClient, logger *slog.Logger) func() {
	type deletionTask struct {
		id         string
		deleteFunc func(ctx context.Context, id string) (string, error)
	}

	deletionTasks := []deletionTask{}

	reportPublished := time.Date(2024, 11, 10, 23, 12, 54, 0, time.UTC)

	label, err := client.CreateLabel(ctx, "id", entity.LabelAddInput{
		Value: "c2-fingerprints",
	})
	if err != nil {
		logger.Error("provisioning label", "error", err)
	}

	deletionTasks = append(deletionTasks, deletionTask{label.ID, client.DeleteLabel})

	// Provision IoC
	observable, err := client.CreateStixCyberObservable(ctx, "id", entity.StixCyberObservableAddInput{
		Type:     graphql.StixCyberObservableTypeIPV4Addr,
		IPv4Addr: graphql.IPv4AddrAddInput{Value: "1.2.3.4"},
	})
	if err != nil {
		logger.Error("provisioning ioc", "error", err)
	}

	iocID := observable.ID

	deletionTasks = append(deletionTasks, deletionTask{observable.ID, client.DeleteStixCyberObservable})

	// Provision report
	report, err := client.CreateReport(ctx, "id", entity.ReportAddInput{
		Name:        "My first report",
		Description: "Hello world!",
		ReportTypes: []graphql.ReportType{graphql.ReportTypeInternal},
		ObjectLabel: []string{"c2-fingerprints"},
		Published:   &reportPublished,
		Objects:     []string{iocID},
	})
	if err != nil {
		logger.Error("provisioning report", "error", err)
	}

	deletionTasks = append(deletionTasks, deletionTask{report.ID, client.DeleteReport})

	return func() {
		for _, task := range deletionTasks {
			if _, err := task.deleteFunc(ctx, task.id); err != nil {
				logger.Error("cleaning up", "error", err)
			}
		}
	}
}

func getLabelIDs(ctx context.Context, client *gocti.OpenCTIAPIClient, labelValues []any) ([]any, error) {
	selectedLabelIDs := []any{}

	// Filter labels: All labels with 'value' in reportLabels
	filterGroup := list.FilterGroup{
		Mode: list.FilterModeAnd,
		Filters: []list.Filter{{
			Mode:     list.FilterModeOr,
			Key:      []string{"value"},
			Operator: list.FilterOperatorEq,
			Values:   labelValues,
		}},
	}

	// List labels with filters
	labels, err := client.ListLabels(ctx, "", false, nil, list.WithFilters(filterGroup))
	if err != nil {
		return nil, fmt.Errorf("unable to list labels: %w", err)
	}

	for _, label := range labels {
		selectedLabelIDs = append(selectedLabelIDs, label.ID)
	}

	return selectedLabelIDs, nil
}

func getReports(
	ctx context.Context,
	client *gocti.OpenCTIAPIClient,
	labelIDs, observableTypes []any,
	cutoffDate time.Time,
) ([]ExtendedReport, error) {
	// Filter reports: All reports with label 'id' in selectedLabelIDs
	filterGroup := list.FilterGroup{
		Mode: list.FilterModeAnd,
		Filters: []list.Filter{
			{
				Mode:     list.FilterModeAnd,
				Key:      []string{"objectLabel"},
				Operator: list.FilterOperatorEq,
				Values:   labelIDs,
			},
		},
	}

	// Filter objects in report:
	// All objects with 'entity_type' in ipTypes AND 'updated_at' later than ipModificationCutoffDate
	obsFilterGroup := list.FilterGroup{
		Mode: list.FilterModeAnd,
		Filters: []list.Filter{{
			Mode:     list.FilterModeOr,
			Key:      []string{"entity_type"},
			Operator: list.FilterOperatorEq,
			Values:   observableTypes,
		}, {
			Mode:     list.FilterModeAnd,
			Key:      []string{"updated_at"},
			Operator: list.FilterOperatorGte,
			Values:   []any{cutoffDate},
		}},
	}

	obsFilterGroupStr, err := obsFilterGroup.ToGraphQL()
	if err != nil {
		return nil, fmt.Errorf("unable to encode filter group to string: %w", err)
	}

	// Completed customAttributes string with objects filter
	completeAttributes := fmt.Sprintf(reportQueryAttributes, obsFilterGroupStr)

	reportBatchSize := 10

	// List reports with filters
	// We cannot use the usual client.ListReports method because we extended the
	// entity.Report type.
	reports, err := api.StructuredList[ExtendedReport, ExtendedReport](ctx, client, completeAttributes, true, nil,
		list.WithFilters(filterGroup),
		list.WithFirst(reportBatchSize),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to list reports: %w", err)
	}

	return reports, nil
}

func getUniqueObservables(reports []ExtendedReport) []Observable {
	output := []Observable{}

	// Extract all unique observables from the reports
	obsIDs := map[string]bool{}

	for _, report := range reports {
		// List all observables, ignore duplicates
		for _, edge := range report.Objects.Edges {
			if _, seen := obsIDs[edge.Node.ID]; !seen {
				obsIDs[edge.Node.ID] = true

				output = append(output, edge.Node)
			}
		}
	}

	return output
}

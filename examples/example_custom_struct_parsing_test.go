// package main_test contains examples of GoCTI usage.
//
// This example shows how to use automatic struct parsing to personalize custom attributes.
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
)

type CustomReport struct {
	ID          string         `gocti:"id"`
	Name        string         `gocti:"name"`
	Description string         `gocti:"description"`
	CreatedBy   CustomIdentity `gocti:"createdBy"`
}

type CustomIdentity struct {
	ID          string   `gocti:"id"`
	Name        string   `gocti:"name"`
	Description string   `gocti:"description"`
	Roles       []string `gocti:"roles"`
}

func Example_customStructParsing() {
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

	callback, id := provisionReport(ctx, client, logger)
	defer callback()

	// Here we can avoid at the same time:
	// - Retrieving unnecessary fields from the server
	// - Rewriting the custom attributes manually to match our struct
	customReportGraphQLAttributes := api.ParseAttributes[CustomReport]()

	// We can use our custom type in a structured query to benefit directly from the decoding
	customReport, err := api.StructuredRead[entity.Report, CustomReport](
		ctx, client, customReportGraphQLAttributes, id,
	)
	if err != nil {
		logger.Error("reading custom report", "error", err)
	}

	fmt.Printf("Retrieved Name: %v\n", customReport.Name)
	fmt.Printf("Created by: %v\n", customReport.CreatedBy.Name)

	// Output:
	// Retrieved Name: Example Custom Report
	// Created by: The Exemplar
}

func provisionReport(ctx context.Context, client *gocti.OpenCTIAPIClient, logger *slog.Logger) (func(), string) {
	now := time.Now()

	identity, err := client.CreateIdentity(ctx, "id", entity.IdentityAddInput{
		Name: "The Exemplar",
		Type: graphql.IdentityTypeIndividual,
	})
	if err != nil {
		logger.Error("provisioning identity", "error", err)
	}

	report, err := client.CreateReport(ctx, "", entity.ReportAddInput{
		Name:        "Example Custom Report",
		Description: "The custom report description",
		Published:   &now,
		CreatedBy:   identity.ID,
	})
	if err != nil {
		logger.Error("provisioning custom report", "error", err)
	}

	return func() {
		if _, err := client.DeleteReport(ctx, report.ID); err != nil {
			logger.Error("cleaning up", "error", err)
		}

		if _, err := client.DeleteIdentity(ctx, identity.ID); err != nil {
			logger.Error("cleaning up", "error", err)
		}
	}, report.ID
}

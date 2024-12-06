// package main_test contains examples of GoCTI usage.
//
// This example shows how large amounts of data can be retrieved from the
// platform using pagination.
package main_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/weisshorn-cyd/gocti"
	"github.com/weisshorn-cyd/gocti/entity"
	"github.com/weisshorn-cyd/gocti/graphql"
	"github.com/weisshorn-cyd/gocti/list"
)

const numEntries = 50

func Example_pagination() {
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

	// Add many entries to the platform
	for i := range numEntries {
		malware, err := client.CreateMalware(ctx, "", entity.MalwareAddInput{
			Name: fmt.Sprintf("Pagination Example Malware %d", i),
		})
		if err != nil {
			logger.Error("creating malware", "error", err)
		}

		// Since this is an example, we clean up after ourselves
		defer func() {
			if _, err := client.DeleteMalware(ctx, malware.ID); err != nil {
				logger.Error("deleting", "malware", malware, "error", err)
			}
		}()
	}

	// List all the entries we just added, 10 at a time
	pageInfo := graphql.PageInfo{}

	malwareNames := []string{}
	pageSize := 10

	nameFilter := list.FilterGroup{
		Mode: list.FilterModeAnd,
		Filters: []list.Filter{{
			Mode:     list.FilterModeOr,
			Key:      []string{"name"},
			Operator: list.FilterOperatorStartsWith,
			Values:   []any{"Pagination Example Malware "},
		}},
	}

	for range numEntries / pageSize {
		malwares, err := client.ListMalwares(ctx, "name", false,
			&pageInfo,                          // provide a PageInfo struct to get information back
			list.WithFirst(pageSize),           // specify how many entries we want per page
			list.WithAfter(pageInfo.EndCursor), // specify cursor position
			list.WithFilters(nameFilter),       // for consistency of testing, filter out everything but our test malwares
		)
		if err != nil {
			logger.Error("listing malware", "error", err)
		}

		fmt.Printf("We retrieved %d malware entities, out of %d\n", len(malwares), pageInfo.GlobalCount)

		for _, malware := range malwares {
			malwareNames = append(malwareNames, malware.Name)
		}
	}

	fmt.Printf("We retrieved %d malware entities in total\n", len(malwareNames))

	// We indeed retrieved all the unique malware entities
	for i, malwareName := range malwareNames {
		expected := fmt.Sprintf("Pagination Example Malware %d", i)
		if malwareName != expected {
			fmt.Printf("Expected to find '%s' but got '%s'\n", expected, malwareName)
		}
	}

	// Output:
	// We retrieved 10 malware entities, out of 50
	// We retrieved 10 malware entities, out of 50
	// We retrieved 10 malware entities, out of 50
	// We retrieved 10 malware entities, out of 50
	// We retrieved 10 malware entities, out of 50
	// We retrieved 50 malware entities in total
}

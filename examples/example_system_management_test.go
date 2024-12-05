// package main_test contains examples of GoCTI usage.
//
// This example shows how GoCTI can be used for system management.
package main_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/weisshorn-cyd/gocti"
	"github.com/weisshorn-cyd/gocti/entity"
	"github.com/weisshorn-cyd/gocti/graphql"
	"github.com/weisshorn-cyd/gocti/list"
	"github.com/weisshorn-cyd/gocti/system"
)

//nolint:cyclop, funlen // ok for an example
func Example_systemManagement() {
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

	// Get ID of the Connector role
	roles, err := client.ListRoles(ctx, "id, name", false, nil, list.WithSearch("connector"))
	if err != nil {
		logger.Error("getting connector role", "error", err)
	}

	if len(roles) != 1 && strings.ToLower(roles[0].Name) != "connector" {
		logger.Error("did not find the expected role", "roles", roles)
		os.Exit(1)
	}

	connectorRole := roles[0]

	// Create group
	group, err := client.CreateGroup(ctx, "id", system.GroupAddInput{
		Name:                 "MyExampleGroup",
		Description:          "A group used in the system management example.",
		DefaultAssignation:   false,
		AutoNewMarking:       true,
		GroupConfidenceLevel: graphql.ConfidenceLevelInput{MaxConfidence: 50},
	})
	if err != nil {
		logger.Error("creating group", "error", err)
	}

	defer func() {
		if _, err := client.DeleteGroup(ctx, group.ID); err != nil {
			logger.Error("deleting", "group", group, "error", err)
		}
	}()

	// Assign the Connector role to the group
	if _, err := group.AssignRole(ctx, client, connectorRole.ID); err != nil {
		logger.Error("assigning role to group", "error", err)
	}

	// Create new user
	user, err := client.CreateUser(ctx, "id", system.UserAddInput{
		Name:      "exampleuser",
		UserEmail: "example@opencti.io",
		Firstname: "Example",
		Lastname:  "McExampleson",
		Password:  "1234",
	})
	if err != nil {
		logger.Error("creating user", "error", err)
	}

	defer func() {
		if _, err := client.DeleteUser(ctx, user.ID); err != nil {
			logger.Error("deleting", "user", user, "error", err)
		}
	}()

	// Assign user to group
	if _, err := user.AssignGroup(ctx, client, group.ID); err != nil {
		logger.Error("assigning user to group", "error", err)
	}

	// Create task
	task, err := client.CreateTask(ctx, "id", entity.TaskAddInput{
		Name:           "Example Task",
		Description:    "This is a task demonstrates the GoCTI system management capabilities.",
		ObjectAssignee: []string{user.ID},
	})
	if err != nil {
		logger.Error("creating task", "error", err)
	}

	defer func() {
		if _, err := client.DeleteTask(ctx, task.ID); err != nil {
			logger.Error("deleting", "task", task, "error", err)
		}
	}()

	// Get Task assigned to a given user
	tasks, err := client.ListTasks(
		ctx,
		"id, description, objectAssignee{id, name}",
		false,
		nil,
		list.WithFilters(
			list.FilterGroup{
				Mode: "and",
				Filters: []list.Filter{
					{
						Key:      []string{"objectAssignee"},
						Values:   []any{user.ID},
						Operator: list.FilterOperatorEq,
						Mode:     list.FilterModeOr,
					},
				},
				FilterGroups: []list.FilterGroup{},
			},
		),
	)
	if err != nil {
		logger.Error("getting tasks", "error", err)
	}

	if len(tasks) > 0 {
		fmt.Printf("Found %d task(s)\n", len(tasks))
		fmt.Printf(
			"The task has the description: '%s'\n",
			tasks[0].Description,
		)
		fmt.Printf(
			"The task is assigned to user: '%s'\n",
			tasks[0].ObjectAssignee[0].Name,
		)
	} else {
		fmt.Printf("Found no task with the given filter\n")
	}

	// Output:
	// Found 1 task(s)
	// The task has the description: 'This is a task demonstrates the GoCTI system management capabilities.'
	// The task is assigned to user: 'exampleuser'
}

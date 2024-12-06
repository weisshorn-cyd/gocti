package system_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	"github.com/weisshorn-cyd/gocti"
	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/graphql"
	"github.com/weisshorn-cyd/gocti/list"
	"github.com/weisshorn-cyd/gocti/system"
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

func TestSystem(t *testing.T) {
	t.Parallel()

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
			name: "Test User",
			args: args{
				createFunc: api.Create[system.User],
				readFunc:   api.Read[system.User],
				listFunc:   api.List[system.User],
				deleteFunc: client.DeleteUser,
				input: system.UserAddInput{
					Name:      "Test User",
					UserEmail: "test@user.com",
					Password:  "test_password",
				},
			},
		}, {
			name: "Test Group",
			args: args{
				createFunc: api.Create[system.Group],
				readFunc:   api.Read[system.Group],
				listFunc:   api.List[system.Group],
				deleteFunc: client.DeleteGroup,
				input: system.GroupAddInput{
					Name: "Test Group",
					GroupConfidenceLevel: graphql.ConfidenceLevelInput{
						MaxConfidence: 100,
						Overrides: []graphql.ConfidenceLevelOverrideInput{
							{EntityType: "report", MaxConfidence: 60},
						},
					},
				},
			},
		}, {
			name: "Test Role",
			args: args{
				createFunc: api.Create[system.Role],
				readFunc:   api.Read[system.Role],
				listFunc:   api.List[system.Role],
				deleteFunc: client.DeleteRole,
				input: system.RoleAddInput{
					Name: "Test Role",
				},
			},
		}, {
			name: "Test StatusTemplate",
			args: args{
				createFunc: api.Create[system.StatusTemplate],
				readFunc:   api.Read[system.StatusTemplate],
				listFunc:   api.List[system.StatusTemplate],
				deleteFunc: client.DeleteStatusTemplate,
				input: system.StatusTemplateAddInput{
					Name:  "TEST_STATUS_TEMPLATE",
					Color: "red",
				},
				orderBy: "name",
			},
		}, {
			name: "Test TaskTemplate",
			args: args{
				createFunc: api.Create[system.TaskTemplate],
				readFunc:   api.Read[system.TaskTemplate],
				listFunc:   api.List[system.TaskTemplate],
				deleteFunc: client.DeleteTaskTemplate,
				input: system.TaskTemplateAddInput{
					Name: "Test TaskTemplate",
				},
			},
		}, {
			name: "Test SubType",
			args: args{
				createFunc: nil,
				readFunc:   api.Read[system.SubType],
				listFunc:   api.List[system.SubType],
				deleteFunc: nil,
				input:      nil,
				orderBy:    "label",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			id := "Status"

			if test.args.createFunc != nil && test.args.deleteFunc != nil {
				id = createEntityAndScheduleDelete(t, client, test.args, test.name)
			}

			if test.args.readFunc != nil {
				readEntity(t, client, test.args, id, test.name)
			}

			if test.args.listFunc != nil {
				listEntities(t, client, test.args, test.name)
			}
		})
	}
}

func TestCaseTemplate(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(t)

	client, err := gocti.NewOpenCTIAPIClient(
		cfg.URL, cfg.Token,
		gocti.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	)
	if err != nil {
		t.Fatalf("cannot create client: %v", err)
	}

	taskID := createEntityAndScheduleDelete(t, client, args{
		createFunc: api.Create[system.TaskTemplate],
		deleteFunc: client.DeleteTaskTemplate,
		input: system.TaskTemplateAddInput{
			Name: "Test TaskTemplate (CaseTemplate)",
		},
	}, "Test TaskTemplate (CaseTemplate)")

	caseArgs := args{
		createFunc: api.Create[system.CaseTemplate],
		readFunc:   api.Read[system.CaseTemplate],
		listFunc:   api.List[system.CaseTemplate],
		deleteFunc: client.DeleteCaseTemplate,
		input: system.CaseTemplateAddInput{
			Name:  "Test CaseTemplate",
			Tasks: []string{taskID},
		},
		orderBy: "name",
	}

	caseID := createEntityAndScheduleDelete(t, client, caseArgs, "Test CaseTemplate")

	readEntity(t, client, caseArgs, caseID, "Test CaseTemplate")
	listEntities(t, client, caseArgs, "Test CaseTemplate")
}

func createEntityAndScheduleDelete(t *testing.T, client *gocti.OpenCTIAPIClient, args args, name string) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create the Entity
	createResp, err := args.createFunc(
		ctx, client, "id", args.input,
	)
	if err != nil {
		t.Fatalf("cannot create %s: %v", name, err)
	}

	idAny, ok := createResp["id"]
	if !ok {
		t.Fatalf("cannot get id after %s creation", name)
	}

	id, ok := idAny.(string)
	if !ok {
		t.Fatalf("cannot get id string after %s creation", name)
	}

	// Setup the Entity's destruction
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_, err := args.deleteFunc(
			ctx, id,
		)
		if err != nil {
			t.Fatalf("failed to remove %s: %v", name, err)
		}
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

	opts := []list.Option{}
	if args.orderBy != "" {
		opts = append(opts, list.WithOrderBy(args.orderBy))
	}

	// Test the List function on the Label
	_, err := args.listFunc(
		ctx, client, "", false, nil, opts...,
	)
	if err != nil {
		t.Errorf("cannot list %s: %v", name, err)
	}
}

func TestWorkbench(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(t)

	client, err := gocti.NewOpenCTIAPIClient(
		cfg.URL, cfg.Token,
		gocti.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	)
	if err != nil {
		t.Fatalf("cannot create client: %v", err)
	}

	type args struct {
		file          api.File
		ErrOnExisting bool
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Empty file",
			args: args{
				file: api.File{
					Name: "test_empty.json",
					Data: []byte{},
					MIME: "application/json",
				},
				ErrOnExisting: false,
			},
		}, {
			name: "Empty file - Double Creation",
			args: args{
				file: api.File{
					Name: "test_empty_double.json",
					Data: []byte{},
					MIME: "application/json",
				},
				ErrOnExisting: true,
			},
		}, {
			name: "Stix Bundle",
			args: args{
				file: api.File{
					Name: "test_bundle.json",
					Data: []byte(`{
  "id": "bundle--4a174b90-fc33-4137-8750-4f4a6e0c3e53",
  "type": "bundle",
  "objects": [
    {
      "name": "Test Threat Actor",
      "first_seen": "2024-11-11T15:49:59+01:00",
      "last_seen": "2024-11-11T15:49:59+01:00",
      "id": "threat-actor--b8aa64c0-3293-400a-886b-3dd5884b8693",
      "type": "threat-actor"
    }
  ]
}`),
					MIME: "application/json",
				},
				ErrOnExisting: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel1 := context.WithTimeout(context.Background(), timeout)
			defer cancel1()

			workbench, err := system.CreateWorkbench(
				ctx, client,
				test.args.file,
				test.args.ErrOnExisting,
			)
			require.NoError(t, err, "creating workbench")

			if test.args.ErrOnExisting {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				_, err := system.CreateWorkbench(ctx, client, test.args.file, test.args.ErrOnExisting)
				require.Error(t, err)
			}

			ctx, cancel2 := context.WithTimeout(context.Background(), timeout)
			defer cancel2()

			_, err = system.DeleteWorkbench(ctx, client, workbench.ID)
			require.NoError(t, err, "deleting workbench")
		})
	}
}

package api_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weisshorn-cyd/gocti"
	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/entity"
	"github.com/weisshorn-cyd/gocti/graphql"
	"github.com/weisshorn-cyd/gocti/list"
)

func TestList(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(t)

	type args struct {
		customAttributes string
		getAll           bool
		pageInfo         *graphql.PageInfo
		opts             []list.Option
	}

	tests := []struct {
		name     string
		listFunc func(
			ctx context.Context,
			client api.Client,
			customAttributes string,
			getAll bool,
			pageInfo *graphql.PageInfo,
			opts ...list.Option,
		) ([]map[string]any, error)
		args         args
		serverResp   string
		wantPageInfo *graphql.PageInfo
		wantErr      bool
	}{
		{
			name:     "api.List ok",
			listFunc: api.List[entity.Label],
			args: args{
				getAll:   false,
				pageInfo: &graphql.PageInfo{},
			},
			serverResp: `
				{"data": {"labels": {
					"edges": [{"node": {"id": "label-id"}}],
					"pageInfo": {
						"startCursor": "start_cursor",
						"endCursor": "end_cursor",
						"hasNextPage": false,
						"hasPreviousPage": false,
						"globalCount": 1
					}
				}}}`,
			wantPageInfo: &graphql.PageInfo{
				StartCursor:     "start_cursor",
				EndCursor:       "end_cursor",
				HasNextPage:     false,
				HasPreviousPage: false,
				GlobalCount:     1,
			},
			wantErr: false,
		}, {
			name:     "api.List with get all",
			listFunc: api.List[entity.Label],
			args: args{
				getAll:   true,
				pageInfo: nil,
			},
			serverResp: `
				{"data": {"labels": {
					"edges": [{"node": {"id": "label-id"}}],
					"pageInfo": {
						"startCursor": "start_cursor",
						"endCursor": "end_cursor",
						"hasNextPage": false,
						"hasPreviousPage": false,
						"globalCount": 3
					}
				}}}`,
			wantPageInfo: nil,
			wantErr:      false,
		}, {
			name:     "api.List with empty returned entity",
			listFunc: api.List[entity.Label],
			args: args{
				getAll:   true,
				pageInfo: nil,
			},
			serverResp: `
			{"data": {"labels": null}}`,
			wantPageInfo: nil,
			wantErr:      false,
		}, {
			name:     "api.List with GraphQL error",
			listFunc: api.List[entity.Label],
			args: args{
				getAll:   false,
				pageInfo: nil,
			},
			serverResp: `
			{"data": {"labels": null},
			"errors":[ {"message": "Mocked Error"}]}`,
			wantPageInfo: nil,
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, err := gocti.NewOpenCTIAPIClient(
				cfg.URL, cfg.Token,
				gocti.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
				gocti.WithTransport(&MockTransport{
					jsonBody: test.serverResp,
				}),
			)
			if err != nil {
				t.Fatalf("cannot create client: %v", err)
			}

			_, err = test.listFunc(
				context.Background(),
				client,
				test.args.customAttributes,
				test.args.getAll,
				test.args.pageInfo,
				test.args.opts...,
			)
			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.wantPageInfo, test.args.pageInfo)
		})
	}
}

func TestStructuredList(t *testing.T) {
	t.Parallel()

	type representative struct {
		Main      string `gocti:"main"`
		Secondary string `gocti:"secondary"`
	}

	type customLabel struct {
		// Different types
		ID        string    `gocti:"id"`
		ToStix    bool      `gocti:"toStix"`
		CreatedAt time.Time `gocti:"created_at"`

		// Nested fields
		Representative representative `gocti:"representative"`

		// Remaining fields
		Other map[string]any `gocti:",remain"`
	}

	cfg := loadConfig(t)

	client, err := gocti.NewOpenCTIAPIClient(
		cfg.URL, cfg.Token,
		gocti.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
		gocti.WithTransport(&MockTransport{
			jsonBody: `
				{"data": {"labels": {
					"edges": [{"node": {
						"id": "label-id",
						"toStix": true,
						"created_at": "2019-03-18T15:46:37.000Z",
						"representative": {
							"main": "label-main",
							"secondary": "label-secondary"
						},
						"unknown_field": "its_value"
					}}]
				}}}`,
		}),
	)
	if err != nil {
		t.Fatalf("cannot create client: %v", err)
	}

	tests := []struct {
		name       string
		test       func() (any, error)
		serverResp string
		want       any
		wantErr    bool
	}{
		{
			name: "Label ok",
			test: func() (any, error) {
				return api.StructuredList[entity.Label, customLabel](
					context.Background(), client, "", false, nil,
				)
			},
			want: []customLabel{{
				ID:        "label-id",
				ToStix:    true,
				CreatedAt: time.Date(2019, time.March, 18, 15, 46, 37, 0, time.UTC),
				Representative: representative{
					Main:      "label-main",
					Secondary: "label-secondary",
				},
				Other: map[string]any{
					"unknown_field": "its_value",
				},
			}},
			wantErr: false,
		}, {
			name: "Field type error",
			test: func() (any, error) {
				return api.StructuredList[entity.Label, struct {
					ID bool `gocti:"id"`
				}](
					context.Background(), client, "", false, nil,
				)
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Non-struct type",
			test: func() (any, error) {
				return api.StructuredList[entity.Label, int](
					context.Background(), client, "", false, nil,
				)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			data, err := test.test()
			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, data)
		})
	}
}

// TestListConfigOptions tests proper hierarchy of list options. Intended priority, from lowest to highest:
// constants in [list.options.go] / env. vars / client option [gocti.Option] / list option [list.Option].
func TestListConfigOptions(t *testing.T) {
	tests := []struct {
		name       string
		inspector  *QueryInspector
		clientOpts []gocti.Option
		listOpts   []list.Option
		envs       map[string]string
		want       map[string]any
		wantErr    bool
	}{
		{
			name:      "base config",
			inspector: &QueryInspector{},
			want: map[string]any{
				"dynamicFrom": nil,
				"dynamicTo":   nil,
				"filters":     nil,
				"first":       1000,
				"orderBy":     "created_at",
				"orderMode":   "asc",
			},
			wantErr: false,
		}, {
			name:      "base config + env",
			inspector: &QueryInspector{},
			envs: map[string]string{
				"GOCTI_PAGE_SIZE":  "50",
				"GOCTI_ORDER_BY":   "updated_at",
				"GOCTI_ORDER_MODE": "desc",
			},
			want: map[string]any{
				"dynamicFrom": nil,
				"dynamicTo":   nil,
				"filters":     nil,
				"first":       50,
				"orderBy":     "updated_at",
				"orderMode":   "desc",
			},
			wantErr: false,
		}, {
			name:      "base config + client option",
			inspector: &QueryInspector{},
			clientOpts: []gocti.Option{
				gocti.WithDefaultPageSize(20),
				gocti.WithDefaultOrderBy("creator"),
				gocti.WithDefaultOrderMode("desc"),
			},
			want: map[string]any{
				"dynamicFrom": nil,
				"dynamicTo":   nil,
				"filters":     nil,
				"first":       20,
				"orderBy":     "creator",
				"orderMode":   "desc",
			},
			wantErr: false,
		}, {
			name:      "base config + list option",
			inspector: &QueryInspector{},
			listOpts: []list.Option{
				list.WithFirst(10),
				list.WithOrderBy("name"),
				list.WithOrderMode(list.OrderModeDesc),
			},
			want: map[string]any{
				"dynamicFrom": nil,
				"dynamicTo":   nil,
				"filters":     nil,
				"first":       10,
				"orderBy":     "name",
				"orderMode":   "desc",
			},
			wantErr: false,
		}, {
			name:      "base config + env + client option",
			inspector: &QueryInspector{},
			envs: map[string]string{
				"GOCTI_PAGE_SIZE":  "50",
				"GOCTI_ORDER_BY":   "updated_at",
				"GOCTI_ORDER_MODE": "desc",
			},
			clientOpts: []gocti.Option{
				gocti.WithDefaultPageSize(20),
				gocti.WithDefaultOrderBy("creator"),
				gocti.WithDefaultOrderMode("asc"),
			},
			want: map[string]any{
				"dynamicFrom": nil,
				"dynamicTo":   nil,
				"filters":     nil,
				"first":       20,
				"orderBy":     "creator",
				"orderMode":   "asc",
			},
			wantErr: false,
		}, {
			name:      "base config + env + client option + list option",
			inspector: &QueryInspector{},
			envs: map[string]string{
				"GOCTI_PAGE_SIZE":  "50",
				"GOCTI_ORDER_BY":   "updated_at",
				"GOCTI_ORDER_MODE": "desc",
			},
			clientOpts: []gocti.Option{
				gocti.WithDefaultPageSize(20),
				gocti.WithDefaultOrderBy("creator"),
				gocti.WithDefaultOrderMode("desc"),
			},
			listOpts: []list.Option{
				list.WithFirst(10),
				list.WithOrderBy("name"),
				list.WithOrderMode(list.OrderModeAsc),
			},
			want: map[string]any{
				"dynamicFrom": nil,
				"dynamicTo":   nil,
				"filters":     nil,
				"first":       10,
				"orderBy":     "name",
				"orderMode":   "asc",
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup env variables
			for k, v := range test.envs {
				t.Setenv(k, v)
			}

			// Create client with additional options
			totalOpts := []gocti.Option{
				gocti.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
				gocti.WithTransport(test.inspector),
			}
			totalOpts = append(totalOpts, test.clientOpts...)

			cfg := loadConfig(t)

			client, err := gocti.NewOpenCTIAPIClient(
				cfg.URL, cfg.Token,
				totalOpts...,
			)
			if err != nil {
				t.Fatalf("cannot create client: %v", err)
			}

			// List reports with additional options
			_, err = api.List[entity.Report](
				context.Background(), client, "", false, nil,
				test.listOpts...,
			)

			// For comparison, as json only uses float64
			for k := range test.want {
				if val, ok := test.want[k].(int); ok {
					test.want[k] = float64(val)
				}
			}

			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, test.inspector.Vars)
		})
	}
}

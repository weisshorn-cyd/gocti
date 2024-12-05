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
)

func TestRead(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(t)

	type args struct {
		customAttributes string
		id               string
	}

	tests := []struct {
		name     string
		readFunc func(
			ctx context.Context,
			client api.Client,
			customAttributes string,
			id string,
		) (map[string]any, error)
		args       args
		serverResp string
		wantErr    bool
	}{
		{
			name:     "api.Read ok",
			readFunc: api.Read[entity.Report],
			args: args{
				id: "report-id",
			},
			serverResp: `{"data": {"report": {"id": "report-id"}}}`,
			wantErr:    false,
		}, {
			name:     "api.Read with empty returned entity",
			readFunc: api.Read[entity.StixCoreRelationship],
			args: args{
				id: "stix-core-relationship-id",
			},
			serverResp: `{"data": {"stixCoreRelationship": null}}`,
			wantErr:    false,
		}, {
			name:     "api.Read with GraphQL error",
			readFunc: api.Read[entity.Label],
			args: args{
				id: "label-id",
			},
			serverResp: `{"data": {"label":null}, "errors":[ {"message": "Mocked Error"} ]}`,
			wantErr:    true,
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

			_, err = test.readFunc(context.Background(), client, test.args.customAttributes, test.args.id)
			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestStructuredRead(t *testing.T) {
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
				{"data": {"label": {
					"id": "label-id",
					"toStix": true,
					"created_at": "2019-03-18T15:46:37.000Z",
					"representative": {
						"main": "label-main",
						"secondary": "label-secondary"
					},
					"unknown_field": "its_value"
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
				return api.StructuredRead[entity.Label, customLabel](
					context.Background(), client, "", "",
				)
			},
			want: customLabel{
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
			},
			wantErr: false,
		}, {
			name: "Field type error",
			test: func() (any, error) {
				return api.StructuredRead[entity.Label, struct {
					ID bool `gocti:"id"`
				}](
					context.Background(), client, "", "",
				)
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Non-struct type",
			test: func() (any, error) {
				return api.StructuredRead[entity.Label, int](
					context.Background(), client, "", "",
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

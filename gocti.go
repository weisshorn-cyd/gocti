package gocti

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/goccy/go-json"
	"github.com/kelseyhightower/envconfig"

	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/list"
	"github.com/weisshorn-cyd/gocti/system"
)

//go:generate go run ./tools/generator.go ./tools

const (
	goctiVersion = "0.37.0"
)

// Config holds all the [OpenCTIAPIClient] parameters that can be set by environment variables.
type Config struct {
	// Api Config
	DefaultTimeout     time.Duration `default:"10s" envconfig:"TIMEOUT"`
	HealthCheckTimeout time.Duration `default:"3s"  envconfig:"HEALTH_CHECK_TIMEOUT"`

	LogLevel slog.Level `default:"info" envconfig:"LOG_LEVEL"`

	// List Config - Overrides 'list/options.go' default values if set
	PageSize  int    `default:"-1"    envconfig:"PAGE_SIZE"`
	OrderBy   string `default:"UNSET" envconfig:"ORDER_BY"`
	OrderMode string `default:"UNSET" envconfig:"ORDER_MODE"`
}

// loadConfig returns the [OpenCTIAPIClient] config as set by environment variables.
func loadConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("GOCTI", &cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to load gocti config: %w", err)
	}

	return &cfg, nil
}

// OpenCTIAPIClient is the main gocti client to interact with the OpenCTI platform.
type OpenCTIAPIClient struct {
	url   string // OpenCTI API url
	token string // OpenCTI API token

	performHealthCheck bool

	config     *Config
	logger     *slog.Logger
	httpClient *http.Client

	impersonating string
}

// NewOpenCTIAPIClient returns a pointer to a properly initialised [OpenCTIAPIClient].
func NewOpenCTIAPIClient(
	url, token string,
	opts ...Option,
) (*OpenCTIAPIClient, error) {
	if url == "" {
		return nil, ErrMissingURL
	}

	if token == "" {
		return nil, ErrMissingToken
	}

	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("bad config: %w", err)
	}

	client := &OpenCTIAPIClient{
		url:        url + "/graphql",
		token:      token,
		config:     cfg,
		logger:     nil,
		httpClient: &http.Client{Timeout: cfg.DefaultTimeout},
	}

	// Apply options to the new client
	for _, opt := range opts {
		opt(client)
	}

	// Create new structured logger if needed
	if client.logger == nil {
		client.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: client.config.LogLevel,
		}))
	}

	// Optional health check (fetching OpenCTI version)
	if client.performHealthCheck {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HealthCheckTimeout)
		defer cancel()

		if err := client.HealthCheck(ctx); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// DefaultPageSize returns the default pagination size for list queries and an ok expression indicating if the value
// has been set (true) or not (false) by the user via env variables or [OpenCTIAPIClient] constructor options.
// The zero value is returned if not set.
func (client *OpenCTIAPIClient) DefaultPageSize() (int, bool) {
	if client.config.PageSize == -1 {
		return 0, false
	}

	return client.config.PageSize, true
}

// DefaultOrderBy returns the default field name for ordering list queries and an ok expression indicating if the value
// has been set (true) or not (false) by the user via env variables or [OpenCTIAPIClient] constructor options.
// The zero value is returned if not set.
func (client *OpenCTIAPIClient) DefaultOrderBy() (string, bool) {
	if client.config.OrderBy == "UNSET" {
		return "", false
	}

	return client.config.OrderBy, true
}

// DefaultOrderMode returns the default mode for ordering list queries and an ok expression indicating if the value
// has been set (true) or not (false) by the user via env variables or [OpenCTIAPIClient] constructor options.
// The zero value is returned if not set.
func (client *OpenCTIAPIClient) DefaultOrderMode() (string, bool) {
	if client.config.OrderMode == "UNSET" {
		return "", false
	}

	return client.config.OrderMode, true
}

// Logger returns the [OpenCTIAPIClient] stored logger.
func (client *OpenCTIAPIClient) Logger() *slog.Logger { return client.logger }

// Impersonate will setup the OpenCTIAPIClient to impersonate the given user for exactly one query.
func (client *OpenCTIAPIClient) Impersonate(ctx context.Context, username string) error {
	users, err := api.StructuredList[system.User, struct {
		ID   string `gocti:"id"`
		Name string `gocti:"name"`
	}](
		ctx, client, "id, name", true, nil,
		list.WithFilters(list.FilterGroup{
			Mode: list.FilterModeAnd,
			Filters: []list.Filter{
				{
					Mode:     list.FilterModeAnd,
					Key:      []string{"name"},
					Operator: list.FilterOperatorEq,
					Values:   []any{username},
				},
			},
		}),
	)
	if err != nil {
		return fmt.Errorf("retrieving user id for %q: %w", username, err)
	}

	if len(users) != 1 {
		return UserNotFoundError{username: username}
	}

	client.impersonating = users[0].ID

	return nil
}

// openCTIGraphQLResponse is the generic response structure for all OpenCTI GraphQL queries.
type openCTIGraphQLResponse struct {
	Data   map[string]any        `json:"data"`
	Errors []OpenCTIGraphQLError `json:"errors"`
}

// mapFileVariables scans the variables from a [Query] in search of entries of type [api.File].
// It returns the list of concerned vars, and a mapping following the pattern:
// `{ [ITERATOR]: ["variables.[PATH_TO_VAR]"] }` e.g. `{ "0": ["variables.files.1"] }`
// As per https://github.com/jaydenseric/graphql-multipart-request-spec
//
// WARNING: [File] variable values in the input map are set to nil ! (As per convention above).
func mapFileVariables(variables map[string]any) ([]api.File, map[int][]string) {
	filesVars := []api.File{}
	mapping := map[int][]string{}

	// Scan variables for `File` types
	for key, value := range variables {
		index := len(mapping)

		if file, ok := value.(api.File); ok {
			// Check for single file variable
			filesVars = append(filesVars, file)
			mapping[index] = []string{"variables." + key}
			variables[key] = nil
		} else if files, ok := value.([]api.File); ok && len(files) > 0 {
			// Check for multiple files variable
			for i, file := range files {
				filesVars = append(filesVars, file)
				mapping[index+i] = []string{fmt.Sprintf("variables.%s.%d", key, i)}
			}

			variables[key] = make([]any, len(files))
		}
	}

	return filesVars, mapping
}

// queryBody writes the appropriate http request body given the query variables mapping.
// It returns the body buffer, and the corresponding content-type.
//
// Standard queries containing no [File] variables must be sent as 'application/json'.
// Whereas if at least 1 [File] variable is present, it must be sent as 'multipart/form-data'.
func queryBody(query string, variables map[string]any) (*bytes.Buffer, string, error) {
	// Scan variables for files
	files, mapping := mapFileVariables(variables)

	// Prepare base GraphQL body
	graphQLBody, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, "", fmt.Errorf("cannot encode standard query into JSON: %w", err)
	}

	if len(files) == 0 {
		// Standard GraphQL query (Content-Type: application/json)
		return bytes.NewBuffer(graphQLBody), "application/json", nil
	}

	return multipartQueryBody(graphQLBody, files, mapping)
}

// mutlipartQueryBody writes an http body request of type 'multipart/form-data'.
// As per https://github.com/jaydenseric/graphql-multipart-request-spec
func multipartQueryBody(graphQLBody []byte, files []api.File, mapping map[int][]string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	defer writer.Close()

	// Part 1: Operations
	ops, err := writer.CreateFormField("operations")
	if err != nil {
		return nil, "", fmt.Errorf("cannot create 'operations' form field: %w", err)
	}

	_, err = ops.Write(graphQLBody)
	if err != nil {
		return nil, "", fmt.Errorf("cannot write 'operations' to body: %w", err)
	}

	// Part 2: Mapping
	maps, err := writer.CreateFormField("map")
	if err != nil {
		return nil, "", fmt.Errorf("cannot create 'map' form field: %w", err)
	}

	jsonMapping, err := json.Marshal(mapping)
	if err != nil {
		return nil, "", fmt.Errorf("cannot encode mapping into JSON: %w", err)
	}

	_, err = maps.Write(jsonMapping)
	if err != nil {
		return nil, "", fmt.Errorf("cannot write 'map' to body: %w", err)
	}

	// Part 3: Files
	for i, file := range files {
		if err := validateFileContent(file); err != nil {
			return nil, "", fmt.Errorf("bad file format: %w", err)
		}

		// Set the appropriate part header
		partHeader := textproto.MIMEHeader{}
		partHeader.Set(
			"Content-Disposition",
			fmt.Sprintf(`form-data; name="%d"; filename="%s"`, i, file.Name),
		)
		partHeader.Set("Content-Type", file.MIME)

		// Create the file part
		filePart, err := writer.CreatePart(partHeader)
		if err != nil {
			return nil, "", fmt.Errorf("cannot create part for file '%s': %w", file.Name, err)
		}

		// Write the file
		_, err = filePart.Write(file.Data)
		if err != nil {
			return nil, "", fmt.Errorf("cannot write file '%s': %w", file.Name, err)
		}
	}

	return body, writer.FormDataContentType(), nil
}

func validateFileContent(file api.File) error {
	if file.Name == "" {
		return ErrEmptyFileName
	}

	if file.MIME == "" {
		return ErrEmptyMIMEType
	}

	return nil
}

// Query sends the provided query string, with optional variables, to
// the OpenCTI GraphQL server and returns the response.
func (client *OpenCTIAPIClient) Query(
	ctx context.Context,
	query string,
	variables map[string]any,
) (map[string]any, error) {
	// Building the Request's Body
	body, contentType, err := queryBody(query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Creating the http request w/ headers
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		client.url,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create the GraphQL request: %w", err)
	}

	// Add the required headers
	client.setOpenCTIHTTPHeaders(&req.Header, contentType)

	// Sending the request to the server
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to contact server: %w", err)
	}
	defer resp.Body.Close()

	// Handling the server response
	if resp.StatusCode != http.StatusOK { // Server error
		return nil, UnexpectedStatusCodeError{resp.Status}
	}

	// Reading the json response
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read server response: %w", err)
	}

	response := openCTIGraphQLResponse{}
	if err = json.Unmarshal(resBody, &response); err != nil {
		return nil, fmt.Errorf("unable to decode GraphQL response: %w", err)
	}

	// Handling GraphQL Errors
	if len(response.Errors) > 0 {
		errs := make([]error, len(response.Errors))
		for i, e := range response.Errors {
			errs[i] = e
		}

		return nil, errors.Join(errs...)
	}

	return response.Data, nil
}

type healthCheckData struct {
	About about `json:"about" mapstructure:"about"`
}

type about struct {
	Version string `json:"version" mapstructure:"version"`
}

// HealthCheck tries to query the OpenCTI version from the server, returning any error it encounters.
func (client *OpenCTIAPIClient) HealthCheck(ctx context.Context) error {
	query := "query {about {version}}"

	response, err := client.Query(ctx, query, nil)
	if err != nil {
		client.logger.Error("Health check failed", "error", err)

		return fmt.Errorf("health check failed: %w", err)
	}

	// Parse response
	healthCheckData := healthCheckData{}
	if err = mapstructure.Decode(response, &healthCheckData); err != nil {
		client.logger.Error("Health check failed", "error", err)

		return fmt.Errorf("unable to parse health check data: %w", err)
	}

	client.logger.Info("OpenCTI health check succeeded", "opencti_version", healthCheckData.About.Version)

	return nil
}

// setOpenCTIHTTPHeaders populates an http request Header with the appropriate headers.
func (client *OpenCTIAPIClient) setOpenCTIHTTPHeaders(header *http.Header, contentType string) {
	header.Set("User-Agent", "gocti/"+goctiVersion)
	header.Set("Content-Type", contentType)
	header.Set("Authorization", "Bearer "+client.token)

	if client.impersonating != "" {
		header.Set("Opencti-Applicant-Id", client.impersonating)
		client.impersonating = ""
	}
}

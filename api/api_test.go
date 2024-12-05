package api_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/goccy/go-json"
	"github.com/kelseyhightower/envconfig"
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

// MockTransport is a RoundTripper that returns a fixed JSON body.
type MockTransport struct {
	jsonBody string
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Create a mock response with the desired JSON body
	jsonBody := []byte(m.jsonBody)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBuffer(jsonBody)),
		Request:    req,
	}

	return resp, nil
}

// QueryInspector is a RoundTripper that inspects the request's body to retrieve a graphql query and its variables.
// It returns a fixed JSON body that mocks a call to the 'reports' GraphQL query.
type QueryInspector struct {
	body  string
	Query string         `json:"query"`
	Vars  map[string]any `json:"variables"`
}

func (inspector *QueryInspector) RoundTrip(request *http.Request) (*http.Response, error) {
	reqBytes, _ := io.ReadAll(request.Body)
	request.Body.Close()
	request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))

	inspector.body = string(reqBytes)
	_ = json.Unmarshal(reqBytes, &inspector)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(`{"data": {"reports": {"edges": []}}}`)),
		Request:    request,
	}

	return resp, nil
}

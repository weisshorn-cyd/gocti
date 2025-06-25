package gocti

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weisshorn-cyd/gocti/api"
)

//nolint:gochecknoglobals // Test variables
var (
	testFile1 = api.File{
		Name: "File1.json",
		Data: []byte{},
		MIME: "application/json",
	}
	testFile2 = api.File{
		Name: "File2.txt",
		Data: []byte{},
		MIME: "text/plain",
	}
	testFileNoName = api.File{
		Name: "",
		Data: []byte{},
		MIME: "application/json",
	}
	testFileNoMIME = api.File{
		Name: "File3.txt",
		Data: []byte{},
		MIME: "",
	}
)

func Test_mapFileVariables(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		vars      map[string]any
		wantInput map[string]any
		wantFiles []api.File
		wantMap   map[int][]string
	}{
		{
			name: "No file",
			vars: map[string]any{
				"string": "string",
				"int":    12,
				"bool":   false,
				"list":   []string{"val1", "val2"},
				"map":    map[any]any{1: "1", "2": 2},
			},
			wantInput: map[string]any{
				"string": "string",
				"int":    12,
				"bool":   false,
				"list":   []string{"val1", "val2"},
				"map":    map[any]any{1: "1", "2": 2},
			},
			wantFiles: []api.File{},
			wantMap:   map[int][]string{},
		}, {
			name: "Single file",
			vars: map[string]any{
				"file": testFile1,
			},
			wantInput: map[string]any{
				"file": nil,
			},
			wantFiles: []api.File{testFile1},
			wantMap:   map[int][]string{0: {"variables.file"}},
		}, {
			name: "Multiple files",
			vars: map[string]any{
				"files": []api.File{
					testFile1,
					testFile2,
				},
			},
			wantInput: map[string]any{
				"files": []any{nil, nil},
			},
			wantFiles: []api.File{testFile1, testFile2},
			wantMap: map[int][]string{
				0: {"variables.files.0"},
				1: {"variables.files.1"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotFiles, gotMap := mapFileVariables(test.vars)
			assert.Len(t, gotMap, len(gotFiles))
			assert.Equal(t, test.wantInput, test.vars, "Input changes mismatch")
			assert.Equal(t, test.wantFiles, gotFiles, "Output files mismatch")
			assert.Equal(t, test.wantMap, gotMap, "Output map mismatch")
		})
	}
}

func Test_queryBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		query    string
		vars     map[string]any
		wantType string
		wantErr  bool
	}{
		{
			name:  "Standard query - no file var",
			query: "query",
			vars: map[string]any{
				"bool":   false,
				"int":    12,
				"string": "v1",
			},
			wantType: "application/json",
			wantErr:  false,
		}, {
			name:  "Upload query - single file",
			query: "upload",
			vars: map[string]any{
				"file": testFile1,
			},
			wantType: "multipart/form-data",
			wantErr:  false,
		}, {
			name:  "Upload query - multiple files",
			query: "upload",
			vars: map[string]any{
				"files": []api.File{
					testFile1,
					testFile2,
				},
			},
			wantType: "multipart/form-data",
			wantErr:  false,
		}, {
			name:  "Upload query - no file name",
			query: "upload",
			vars: map[string]any{
				"file": testFileNoName,
			},
			wantType: "multipart/form-data",
			wantErr:  true,
		}, {
			name:  "Upload query - no MIME type",
			query: "upload",
			vars: map[string]any{
				"file": testFileNoMIME,
			},
			wantType: "multipart/form-data",
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, gotType, err := queryBody(test.query, test.vars)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, gotType, test.wantType)
			}
		})
	}
}

// MockTransport is a RoundTripper that returns a fixed JSON body.
type MockTransport struct {
	header   http.Header
	jsonBody string
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.header = req.Header.Clone()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(m.jsonBody)),
		Request:    req,
	}

	return resp, nil
}

func TestOpenCTIAPIClient_Impersonate(t *testing.T) {
	t.Parallel()

	mock := &MockTransport{jsonBody: `{"data":{"users":{"edges":[{"node":{"id":"test-user-id"}}]}}}`}
	client, err := NewOpenCTIAPIClient("url", "token", WithTransport(mock))
	require.NoError(t, err)

	err = client.Impersonate(t.Context(), "test-user")
	require.NoError(t, err)

	_, err = client.ListUsers(t.Context(), "id", true, nil)
	require.NoError(t, err)

	assert.Equal(t, "test-user-id", mock.header.Get("Opencti-Applicant-Id"))

	_, err = client.ListUsers(t.Context(), "id", true, nil)
	require.NoError(t, err)

	assert.Empty(t, mock.header.Get("Opencti-Applicant-Id"))
}

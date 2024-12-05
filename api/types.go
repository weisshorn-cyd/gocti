package api

import (
	"context"
	"log/slog"
)

// Defaults defines the interface to access various configured default values.
type Defaults interface {
	// DefaultPageSize returns the default pagination size for list queries and an ok expression indicating if the value
	// has been set (true) or not (false) by the user via env variables or [OpenCTIAPIClient] constructor options.
	// The zero value is returned if not set.
	DefaultPageSize() (int, bool)

	// DefaultOrderBy returns the default field name for ordering list queries and an ok expression indicating if the value
	// has been set (true) or not (false) by the user via env variables or [OpenCTIAPIClient] constructor options.
	// The zero value is returned if not set.
	DefaultOrderBy() (string, bool)

	// DefaultOrderMode returns the default mode for ordering list queries and an ok expression indicating if the value
	// has been set (true) or not (false) by the user via env variables or [OpenCTIAPIClient] constructor options.
	// The zero value is returned if not set.
	DefaultOrderMode() (string, bool)
}

type Client interface {
	Query(ctx context.Context, query string, variables map[string]any) (map[string]any, error)
	Logger() *slog.Logger
	Defaults
}

type EditInput struct {
	Key        string        `json:"key,omitempty"`
	ObjectPath string        `json:"object_path,omitempty"`
	Value      any           `json:"value,omitempty"`
	Operation  EditOperation `json:"operation,omitempty"`
}

type EditOperation string

const (
	EditOperationAdd     EditOperation = "add"
	EditOperationReplace EditOperation = "replace"
	EditOperationRemove  EditOperation = "remove"
)

type File struct {
	Name string
	Data []byte
	MIME string
}

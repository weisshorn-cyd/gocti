package gocti

import (
	"errors"
	"fmt"
)

var (
	ErrMissingURL   = errors.New("an URL must be set")
	ErrMissingToken = errors.New("a token must be set")

	ErrEmptyFileName = errors.New("file has empty filename")
	ErrEmptyMIMEType = errors.New("file has empty MIME type")
)

type UnexpectedStatusCodeError struct{ status string }

func (err UnexpectedStatusCodeError) Error() string {
	return "unexpected http status code: " + err.status
}

type OpenCTIGraphQLError struct {
	Name       string                  `json:"name"`
	Message    string                  `json:"message"`
	TimeThrown string                  `json:"time_thrown"`
	Data       OpenCTIGraphQLErrorData `json:"data"`
}

type OpenCTIGraphQLErrorData struct {
	HTTPStatus int    `json:"http_status"`
	Genre      string `json:"genre"`
	Reason     string `json:"reason"`
}

func (err OpenCTIGraphQLError) Error() string {
	if err.Name == "" {
		err.Name = "Message"
	}

	if err.Data.Reason != "" {
		err.Message = err.Data.Reason
	}

	return fmt.Sprintf("GraphQL error: %s: %s", err.Name, err.Message)
}

type UserNotFoundError struct {
	username string
}

func (err UserNotFoundError) Error() string {
	return fmt.Sprintf("could not find user with name %q", err.username)
}

package gocti

import (
	"log/slog"
	"net/http"
	"time"
)

type Option func(c *OpenCTIAPIClient)

func WithTransport(transport http.RoundTripper) Option {
	return func(c *OpenCTIAPIClient) {
		c.httpClient.Transport = transport
	}
}

func WithHealthCheck() Option {
	return func(c *OpenCTIAPIClient) {
		c.performHealthCheck = true
	}
}

func WithHealthCheckTimeout(timeout time.Duration) Option {
	return func(c *OpenCTIAPIClient) {
		c.config.HealthCheckTimeout = timeout
	}
}

func WithLogLevel(level slog.Level) Option {
	return func(c *OpenCTIAPIClient) {
		c.config.LogLevel = level
	}
}

func WithLogger(log *slog.Logger) Option {
	return func(c *OpenCTIAPIClient) {
		c.logger = log
	}
}

func WithDefaultTimeout(timeout time.Duration) Option {
	return func(c *OpenCTIAPIClient) {
		c.config.DefaultTimeout = timeout
	}
}

func WithDefaultPageSize(pageSize int) Option {
	return func(c *OpenCTIAPIClient) {
		c.config.PageSize = pageSize
	}
}

func WithDefaultOrderBy(orderBy string) Option {
	return func(c *OpenCTIAPIClient) {
		c.config.OrderBy = orderBy
	}
}

func WithDefaultOrderMode(orderMode string) Option {
	return func(c *OpenCTIAPIClient) {
		c.config.OrderMode = orderMode
	}
}

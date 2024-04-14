// Package logger contains objects and methods for logging the events.
package logger

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type (
	// ResponseData contains response data.
	ResponseData struct {
		Body   *bytes.Buffer
		Status int
		Size   int
	}

	// LoggingResponseWriter is an implementation of http.ResponseWriter.
	LoggingResponseWriter struct {
		http.ResponseWriter
		ResponseData *ResponseData
	}
)

// Log is singleton of events logger.
var Log *zap.Logger = zap.NewNop()

// Init initializes logger singleton with the appropriate atomic level.
func Init(ctx context.Context, level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("Init: parse level errors %w", err)
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("Init: logger build error %w", err)
	}
	Log = zl
	return nil
}

// WriteHeader implements writing the header and status code capturing.
func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode
}

// Write implements writing the response and capturing the body size
// and body itself.
func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return size, fmt.Errorf("Write: response write %w", err)
	}
	r.ResponseData.Size += size
	r.ResponseData.Body.Write(b)
	return size, nil
}

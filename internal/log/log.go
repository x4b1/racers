package log

import (
	"context"
)

type Fields map[string]interface{}

type Logger interface {
	Debug(ctx context.Context, msg string, fields Fields)
	Error(ctx context.Context, err error, fields Fields)
	Info(ctx context.Context, msg string, fields Fields)
}

type Noop struct{}

func (Noop) Debug(ctx context.Context, msg string, fields Fields) {}
func (Noop) Error(ctx context.Context, err error, fields Fields)  {}
func (Noop) Info(ctx context.Context, msg string, fields Fields)  {}

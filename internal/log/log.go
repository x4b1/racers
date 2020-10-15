package log

import (
	"context"
)

type Fields map[string]interface{}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, err error, args ...interface{})
	Info(ctx context.Context, args ...interface{})
}

type Noop struct{}

func (Noop) Debug(ctx context.Context, msg string, args ...interface{}) {}
func (Noop) Error(ctx context.Context, err error, args ...interface{})  {}
func (Noop) Info(ctx context.Context, args ...interface{})              {}

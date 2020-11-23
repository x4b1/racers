package log

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

type Payload map[string]interface{}

type Logger interface {
	Debug(ctx context.Context, msg string, payload Payload)
	Error(ctx context.Context, err error, payload Payload)
	Info(ctx context.Context, msg string, payload Payload)
}

func NewLogrus(out io.Writer) (LogrusLogger, error) {
	l := logrus.New()

	l.SetOutput(out)
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetLevel(logrus.DebugLevel)

	return LogrusLogger{l}, nil
}

var _ Logger = LogrusLogger{}

type LogrusLogger struct {
	logger *logrus.Logger
}

func (l LogrusLogger) Debug(ctx context.Context, msg string, payload Payload) {
	l.logger.WithContext(ctx).WithFields(logrus.Fields(payload)).Debug(msg)
}

func (l LogrusLogger) Error(ctx context.Context, err error, payload Payload) {
	l.logger.WithContext(ctx).WithFields(logrus.Fields(payload)).WithError(err).Error(err)
}

func (l LogrusLogger) Info(ctx context.Context, msg string, payload Payload) {
	l.logger.WithContext(ctx).WithFields(logrus.Fields(payload)).Info(msg)
}

var _ Logger = NoopLogger{}

type NoopLogger struct{}

func (NoopLogger) Debug(ctx context.Context, msg string, payload Payload) {}
func (NoopLogger) Error(ctx context.Context, err error, payload Payload)  {}
func (NoopLogger) Info(ctx context.Context, msg string, payload Payload)  {}

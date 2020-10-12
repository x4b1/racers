package log

import (
	"context"
	"io"

	"github.com/xabi93/racers/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(env config.Env, out io.Writer) (ZapLogger, error) {
	var (
		zapLog *zap.Logger
		err    error
	)

	switch env {
	case config.Local:
		zapLog, err = zap.NewDevelopment()
	case config.Prod:
		zapLog, err = zap.NewProduction()
	}
	if err != nil {
		return ZapLogger{}, err
	}

	zapLog = zapLog.WithOptions(zap.ErrorOutput(zapcore.AddSync(out)))

	return ZapLogger{zapLog}, nil
}

var _ Logger = ZapLogger{}

type ZapLogger struct {
	logger *zap.Logger
}

func (l ZapLogger) Debug(ctx context.Context, msg string, fields Fields) {
	l.logger.Debug(msg, l.fields(fields)...)
}

func (l ZapLogger) Error(ctx context.Context, err error, fields Fields) {
	l.logger.Error(err.Error(), append(l.fields(fields), zap.Error(err))...)
}

func (l ZapLogger) Info(ctx context.Context, msg string, fields Fields) {
	l.logger.Info(msg, l.fields(fields)...)
}

func (ZapLogger) fields(fields Fields) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return zapFields
}

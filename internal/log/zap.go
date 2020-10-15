package log

import (
	"context"
	"io"

	"github.com/xabi93/racers/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(e config.Env, out io.Writer) (ZapLogger, error) {
	var (
		zapLog *zap.Logger
		err    error
	)

	switch e {
	case config.Local:
		zapLog, err = zap.NewDevelopment()
	case config.Prod:
		zapLog, err = zap.NewProduction()
	}
	if err != nil {
		return ZapLogger{}, err
	}

	return ZapLogger{zapLog.WithOptions(zap.ErrorOutput(zapcore.AddSync(out))).Sugar()}, nil
}

var _ Logger = ZapLogger{}

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func (l ZapLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.logger.With(ctx).Debugw(msg, args...)
}

func (l ZapLogger) Error(ctx context.Context, err error, args ...interface{}) {
	l.logger.With(ctx).Error([]interface{}{zap.Error(err), args}...)
}

func (l ZapLogger) Info(ctx context.Context, args ...interface{}) {
	l.logger.With(ctx).Info(args...)
}

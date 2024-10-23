package clog

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger field keys
const KeyClientId = "client_id"

type key string

var loggerKey key = "logger"
var logger *zap.Logger
var logLevel zap.AtomicLevel

func init() {
	config := zap.NewProductionConfig()
	logLevel = config.Level
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger = zap.Must(config.Build())
}

// SetLogLevel attempts to set the log level globally to the given value.
// If an error is encountered, it sets the global log level to zapcore.InfoLevel
func SetLogLevel(l string) {
	if err := logLevel.UnmarshalText([]byte(l)); err != nil {
		logLevel.SetLevel(zapcore.InfoLevel)
	}
}

func NewContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func NewContextWithDefaultLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func MustFromContext(ctx context.Context) *zap.Logger {
	u, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		// TODO: Graceful shutdown instead of panic
		panic("logger not found in context")
	}
	return u
}

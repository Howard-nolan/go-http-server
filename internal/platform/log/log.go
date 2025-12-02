package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps a zap logger so callers can keep using familiar sugar helpers.
type Logger struct {
	*zap.SugaredLogger
	base *zap.Logger
}

// New returns a production-ready zap logger.
func New() *Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	base, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		SugaredLogger: base.Sugar(),
		base:          base,
	}
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() {
	_ = l.base.Sync()
}

// Desugar exposes the underlying structured logger.
func (l *Logger) Desugar() *zap.Logger {
	return l.base
}

package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zap field is same as key value pairs that is attached to logs
type Field = zap.Field 

type ZapLogger struct {
	l *zap.Logger 
}

// This is just a trick to make a unique key for context
var ctxKey = &struct{}{}

func NewZapLogger(level string, json bool) (Logger, error){
	var cfg zap.Config

	// json means your logs will be in json format
	if json {
		cfg = zap.NewProductionConfig()
	}else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
	}

	// you can think of zapcore.Level as zap enum level
	var lvl zapcore.Level 

	// convert string level to zapcore.Level
	// For example: "info" -> zapcore.InfoLevel
	if err := lvl.UnmarshalText([]byte(level)); err == nil {
		// set the log level to the log configuration
		// Atomic level means you can change the log level at runtime
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	}

	zl, err := cfg.Build()

	if err != nil {
		return nil, err 
	}

	return &ZapLogger{l: zl}, nil
}

func (z *ZapLogger) Info(msg string, fields ...Field){
	z.l.Info(msg, fields...)
}

func (z *ZapLogger) Debug(msg string, fields ...Field){
	z.l.Debug(msg, fields...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field){
	z.l.Warn(msg, fields...)
}	

func (z *ZapLogger) Error(msg string, fields ...Field){
	z.l.Error(msg, fields...)
}

func (z *ZapLogger) Fatal(msg string, fields ...Field){
	z.l.Fatal(msg, fields...)
}

// Add fields that automatically attached to logs to all of future logs
func (z *ZapLogger) With(fields ...Field) Logger{
	return &ZapLogger{l: z.l.With(fields...)}
}

// Flush any buffered log entries to output file (file /stdout / etc)
func (z *ZapLogger) Sync() error{return z.l.Sync()}

// get the zap logger from context
func (z *ZapLogger) FromContext(ctx context.Context) Logger{
	if v := ctx.Value(ctxKey); v != nil {
		if l, ok := v.(*zap.Logger); ok {
			return &ZapLogger{l: l}
		}
	}

	return z 
}

// insert or set zap logger into context
func ContextWithLogger(ctx context.Context, logger Logger) context.Context{
	if zl, ok := logger.(*ZapLogger); ok {
		return context.WithValue(ctx, ctxKey, zl.l)
	}

	return ctx 
}
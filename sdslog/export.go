package sdslog

import (
	"context"
	"log/slog"
)

type (
	Attr  = slog.Attr
	Level = slog.Level
)

var (
	String   = slog.String
	Int64    = slog.Int64
	Int      = slog.Int
	Uint64   = slog.Uint64
	Float64  = slog.Float64
	Bool     = slog.Bool
	Time     = slog.Time
	Duration = slog.Duration
)

func With(args ...any) Logger {
	return L(nil).With(args...)
}

func WithError(err any) Logger {
	return L(nil).WithError(err)
}

func WithAttrs(attrs map[string]any) Logger {
	return L(nil).WithAttrs(attrs)
}

func WithAttr(k string, v any) Logger {
	return L(nil).WithAttr(k, v)
}

func Log(ctx context.Context, level Level, msg string, args ...any) {
	L(nil).Log(ctx, level, msg, args...)
}

func Debug(msg string, args ...any) {
	L(nil).Debug(msg, args...)
}

func Info(msg string, args ...any) {
	L(nil).Info(msg, args...)
}

func Warn(msg string, args ...any) {
	L(nil).Warn(msg, args...)
}

func Error(msg string, args ...any) {
	L(nil).Error(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	L(nil).DebugContext(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	L(nil).InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	L(nil).WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	L(nil).ErrorContext(ctx, msg, args...)
}

func Debugf(msg string, a ...any) {
	L(nil).Debugf(msg, a...)
}

func Infof(msg string, a ...any) {
	L(nil).Infof(msg, a...)
}

func Warnf(msg string, a ...any) {
	L(nil).Warnf(msg, a...)
}

func Errorf(msg string, a ...any) {
	L(nil).Errorf(msg, a...)
}

func DebugfContext(ctx context.Context, msg string, a ...any) {
	L(nil).DebugfContext(ctx, msg, a...)
}

func InfofContext(ctx context.Context, msg string, a ...any) {
	L(nil).InfofContext(ctx, msg, a...)
}

func WarnfContext(ctx context.Context, msg string, a ...any) {
	L(nil).WarnfContext(ctx, msg, a...)
}

func ErrorfContext(ctx context.Context, msg string, a ...any) {
	L(nil).ErrorfContext(ctx, msg, a...)
}

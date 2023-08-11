package sdslog

import (
	"context"
	"fmt"
	"log/slog"
)

// format

func Debugf(msg string, args ...any) {
	slog.Debug(fmt.Sprintf(msg, args...))
}

func Infof(msg string, args ...any) {
	slog.Info(fmt.Sprintf(msg, args...))
}

func Warnf(msg string, args ...any) {
	slog.Warn(fmt.Sprintf(msg, args...))
}

func Errorf(msg string, args ...any) {
	slog.Error(fmt.Sprintf(msg, args...))
}

func DebugfContext(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, fmt.Sprintf(msg, args...))
}

func InfofContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, fmt.Sprintf(msg, args...))
}

func WarnfContext(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, fmt.Sprintf(msg, args...))
}

func ErrorfContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

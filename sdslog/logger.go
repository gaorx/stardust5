package sdslog

import (
	"context"
	"fmt"
	"github.com/gaorx/stardust5/sderr"
	"log/slog"
)

type Logger struct {
	*slog.Logger
}

func L(l *slog.Logger) Logger {
	if l == nil {
		l = slog.Default()
	}
	return Logger{l}
}

func (l Logger) With(args ...any) Logger {
	return Logger{l.Logger.With(args...)}
}

func (l Logger) WithError(err any) Logger {
	if err == nil {
		return l
	}
	return Logger{l.Logger.With("error", sderr.AsErr(err).Error())}
}

func (l Logger) WithAttrs(attrs map[string]any) Logger {
	if len(attrs) <= 0 {
		return l
	}
	return Logger{l.Logger.With(MapToArgs(attrs)...)}
}

func (l Logger) WithAttr(k string, v any) Logger {
	return Logger{l.Logger.With(k, v)}
}

func (l Logger) WithCall(call string) Logger {
	if call == "" {
		return l
	}
	return l.WithAttr("call", call)
}

func (l Logger) WithFunc(f string) Logger {
	if f == "" {
		return l
	}
	return l.WithAttr("func", f)
}

func (l Logger) WithAPI(api string) Logger {
	if api == "" {
		return l
	}
	return l.WithAttr("api", api)
}

func (l Logger) Debugf(format string, a ...any) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l Logger) Infof(format string, a ...any) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l Logger) Warnf(format string, a ...any) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l Logger) Errorf(format string, a ...any) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l Logger) DebugfContext(ctx context.Context, format string, a ...any) {
	l.DebugContext(ctx, fmt.Sprintf(format, a...))
}

func (l Logger) InfofContext(ctx context.Context, format string, a ...any) {
	l.InfoContext(ctx, fmt.Sprintf(format, a...))
}

func (l Logger) WarnfContext(ctx context.Context, format string, a ...any) {
	l.WarnContext(ctx, fmt.Sprintf(format, a...))
}

func (l Logger) ErrorfContext(ctx context.Context, format string, a ...any) {
	l.ErrorContext(ctx, fmt.Sprintf(format, a...))
}

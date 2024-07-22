package sdbun

import (
	"context"
	"github.com/gaorx/stardust5/sdslog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"time"
)

func LoggerOf(name string) bun.QueryHook {
	switch name {
	case "", "discard", "disable":
		return discardLogger{}
	case "default", "bun":
		return bundebug.NewQueryHook(bundebug.WithVerbose(true), bundebug.FromEnv("BUNDEBUG"))
	case "builtin", "stardust", "sd", "slog":
		return builtinLogger{}
	default:
		return discardLogger{}
	}
}

type discardLogger struct{}

func (h discardLogger) BeforeQuery(ctx context.Context, e *bun.QueryEvent) context.Context {
	return ctx
}

func (h discardLogger) AfterQuery(ctx context.Context, e *bun.QueryEvent) {
}

type builtinLogger struct{}

func (h builtinLogger) BeforeQuery(ctx context.Context, e *bun.QueryEvent) context.Context {
	return ctx
}

func (h builtinLogger) AfterQuery(ctx context.Context, e *bun.QueryEvent) {
	elapsed := time.Now().Sub(e.StartTime).String()
	if e.Err != nil {
		sdslog.WithError(e.Err).
			WithAttr("q", e.Query).
			WithAttr("elapsed", elapsed).
			WithAttr("op", e.Operation()).
			Error("query error")
	} else {
		sdslog.WithAttr("q", e.Query).
			WithAttr("elapsed", elapsed).
			WithAttr("op", e.Operation()).
			Debug("query done")
	}
}

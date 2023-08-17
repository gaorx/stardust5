package sdgorm

import (
	"context"
	"errors"
	"github.com/gaorx/stardust5/sdslog"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gaorx/stardust5/sdtime"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	gormutils "gorm.io/gorm/utils"
)

var (
	LoggerDiscard        = gormlogger.Discard
	LoggerGormDefault    = gormlogger.Default
	LoggerPlainInfo      = newGormLogger(gormlogger.Info, false)
	LoggerPlainWarn      = newGormLogger(gormlogger.Warn, false)
	LoggerPlainError     = newGormLogger(gormlogger.Error, false)
	LoggerPlainSilent    = newGormLogger(gormlogger.Silent, false)
	LoggerColorfulInfo   = newGormLogger(gormlogger.Info, true)
	LoggerColorfulWarn   = newGormLogger(gormlogger.Warn, true)
	LoggerColorfulError  = newGormLogger(gormlogger.Error, true)
	LoggerColorfulSilent = newGormLogger(gormlogger.Silent, true)
	LoggerBuiltin        = newBuiltinLogger(200 * time.Millisecond)
)

func LoggerOf(name string) gormlogger.Interface {
	switch strings.ToLower(name) {
	case "", "default":
		return LoggerDiscard
	case "gorm_default":
		return LoggerGormDefault
	case "discard", "disable":
		return LoggerDiscard
	case "builtin", "stardust", "sd", "slog":
		return LoggerBuiltin
	case "plain_info", "info":
		return LoggerPlainInfo
	case "plain_warn", "warn":
		return LoggerPlainWarn
	case "plain_error", "error":
		return LoggerPlainError
	case "plain_silent", "silent":
		return LoggerPlainSilent
	case "colorful_info", "color_info":
		return LoggerColorfulInfo
	case "colorful_warn", "color_warn":
		return LoggerColorfulInfo
	case "colorful_error", "color_error":
		return LoggerColorfulError
	case "colorful_silent", "color_silent":
		return LoggerColorfulSilent
	default:
		return LoggerDiscard
	}
}

func newGormLogger(level gormlogger.LogLevel, colorful bool) gormlogger.Interface {
	return gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: false,
			Colorful:                  colorful,
		},
	)
}

type builtinLogger struct {
	slowThreshold         time.Duration
	skipErrRecordNotFound bool
}

func newBuiltinLogger(slowThreshold time.Duration) *builtinLogger {
	return &builtinLogger{
		slowThreshold:         slowThreshold,
		skipErrRecordNotFound: true,
	}
}

func (l *builtinLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *builtinLogger) Info(ctx context.Context, s string, args ...interface{}) {
	sdslog.InfofContext(ctx, s, args...)
}

func (l *builtinLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	sdslog.WarnfContext(ctx, s, args...)
}

func (l *builtinLogger) Error(ctx context.Context, s string, args ...interface{}) {
	sdslog.ErrorfContext(ctx, s, args...)
}

func (l *builtinLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	attrs := []any{
		sdslog.Float64("elapsed", sdtime.ToMillisF(elapsed)),
		sdslog.String("line", gormutils.FileWithLineNum()),
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.skipErrRecordNotFound) {
		sdslog.With(attrs...).WithError(err).ErrorContext(ctx, sql)
		return
	}

	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		sdslog.With(attrs...).WarnContext(ctx, sql)
		return
	}

	sdslog.With(attrs...).DebugContext(ctx, sql)
}

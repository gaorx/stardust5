package sdslog

import (
	"fmt"
	"github.com/samber/lo"
	"log/slog"
	"time"
)

func MapToArgs(m map[string]any) []any {
	return lo.ToAnySlice(MapToAttrs(m))
}

func MapToAttrs(m map[string]any) []slog.Attr {
	var attrs []slog.Attr
	for k, v := range m {
		switch v1 := v.(type) {
		case nil:
			attrs = append(attrs, slog.String(k, "<nil>"))
		case string:
			attrs = append(attrs, slog.String(k, v1))
		case int:
			attrs = append(attrs, slog.Int(k, v1))
		case bool:
			attrs = append(attrs, slog.Bool(k, v1))
		case int64:
			attrs = append(attrs, slog.Int64(k, v1))
		case uint64:
			attrs = append(attrs, slog.Uint64(k, v1))
		case float64:
			attrs = append(attrs, slog.Float64(k, v1))
		case time.Time:
			attrs = append(attrs, slog.Time(k, v1))
		case time.Duration:
			attrs = append(attrs, slog.Duration(k, v1))
		case error:
			if v1 != nil {
				attrs = append(attrs, slog.String(k, v1.Error()))
			} else {
				attrs = append(attrs, slog.String(k, ""))
			}
		default:
			attrs = append(attrs, slog.String(k, fmt.Sprintf("%v", v)))
		}
	}
	return attrs
}

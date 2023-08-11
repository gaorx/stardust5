package sdslog

import (
	"log/slog"
)

func WithError(err error) *slog.Logger {
	if err != nil {
		return slog.With("error", err.Error())
	} else {
		return slog.Default()
	}
}

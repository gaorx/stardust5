package sdslog

import (
	"log/slog"
)

func Setup(opts *Options) {
	slog.SetDefault(Must(opts))
}

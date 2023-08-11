package sdslog

import (
	"github.com/samber/lo"
	"log/slog"
)

func Setup(opts *Options) {
	slog.SetDefault(lo.Must(New(opts)))
}

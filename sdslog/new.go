package sdslog

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
)

type Format string

type Options struct {
	Level     string   `json:"level" toml:"level" yaml:"level"`
	Format    string   `json:"format" toml:"format" yaml:"format"`
	Outputs   []string `json:"output" toml:"output" yaml:"output"`
	AddSource bool     `json:"add_source" toml:"add_source" yaml:"add_source"`
}

func New(opts *Options) (*slog.Logger, error) {
	opts1 := lo.FromPtr(opts)

	w, err := newWriter(opts1.Outputs)
	if err != nil {
		return nil, sderr.WithStack(err)
	}

	// slog options
	slogOpts := &slog.HandlerOptions{AddSource: opts1.AddSource}

	// level
	switch strings.ToLower(opts1.Level) {
	case "debug":
		slogOpts.Level = slog.LevelDebug
	case "info", "":
		slogOpts.Level = slog.LevelInfo
	case "warn":
		slogOpts.Level = slog.LevelWarn
	case "error":
		slogOpts.Level = slog.LevelError
	default:
		return nil, sderr.NewWith("illegal level", opts1.Level)
	}

	// handler
	var h slog.Handler
	switch opts1.Format {
	case "", "text":
		h = slog.NewTextHandler(w, slogOpts)
	case "json":
		h = slog.NewJSONHandler(w, slogOpts)
	default:
		return nil, sderr.NewWith("illegal format", opts1.Format)
	}

	// go
	return slog.New(h), nil
}

func newWriter(outputs []string) (io.Writer, error) {
	one := func(output string) (io.Writer, error) {
		if output == "" || strings.ToLower(output) == "stdout" {
			return os.Stdout, nil
		} else if strings.ToLower(output) == "stderr" {
			return os.Stderr, nil
		} else if strings.ToLower(output) == "discard" {
			return io.Discard, nil
		} else {
			w, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return nil, sderr.WrapWith(err, "open file error", output)
			}
			return w, nil
		}
	}
	if len(outputs) <= 0 {
		return one("stdout")
	} else if len(outputs) == 1 {
		return one(outputs[0])
	} else {
		var writers []io.Writer
		for _, output := range outputs {
			w, err := one(output)
			if err != nil {
				return nil, sderr.WithStack(err)
			}
			writers = append(writers, w)
		}
		return io.MultiWriter(writers...), nil
	}
}

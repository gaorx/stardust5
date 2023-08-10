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

const (
	Text Format = "text"
	Json Format = "json"
)

type Options struct {
	Format    Format   `json:"format" toml:"format" yaml:"format"`
	Outputs   []string `json:"output" toml:"output" yaml:"output"`
	AddSource bool     `json:"add_source" toml:"add_source" yaml:"add_source"`
}

func Must(opts *Options) *slog.Logger {
	return lo.Must(New(opts))
}

func New(opts *Options) (*slog.Logger, error) {
	opts1 := lo.FromPtr(opts)

	w, err := newWriter(opts1.Outputs)
	if err != nil {
		return nil, sderr.WithStack(err)
	}

	slogOpts := &slog.HandlerOptions{AddSource: opts1.AddSource}
	switch opts1.Format {
	case "", Text:
		return slog.New(slog.NewTextHandler(w, slogOpts)), nil
	case Json:
		return slog.New(slog.NewJSONHandler(w, slogOpts)), nil
	default:
		return nil, sderr.NewWith("illegal format", opts1.Format)
	}
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

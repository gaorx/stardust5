package sdbackoff

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
)

type (
	BackOff = backoff.BackOff
	Ticker  = backoff.Ticker
)

type ExponentialOptions struct {
	InitialInterval     time.Duration
	RandomizationFactor float64
	Multiplier          float64
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration
}

func Exponential(opts *ExponentialOptions) BackOff {
	opts1 := lo.FromPtr(opts)
	if opts1.InitialInterval <= 0 {
		opts1.InitialInterval = backoff.DefaultInitialInterval
	}
	if opts1.RandomizationFactor <= 0.0 {
		opts1.RandomizationFactor = backoff.DefaultRandomizationFactor
	}
	if opts1.Multiplier <= 0.0 {
		opts1.Multiplier = backoff.DefaultMultiplier
	}
	if opts1.MaxInterval <= 0 {
		opts1.MaxInterval = backoff.DefaultMaxInterval
	}
	if opts1.MaxElapsedTime <= 0 {
		opts1.MaxElapsedTime = backoff.DefaultMaxElapsedTime
	}
	b := &backoff.ExponentialBackOff{
		InitialInterval:     opts1.InitialInterval,
		RandomizationFactor: opts1.RandomizationFactor,
		Multiplier:          opts1.Multiplier,
		MaxInterval:         opts1.MaxInterval,
		MaxElapsedTime:      opts1.MaxElapsedTime,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	return b
}

func Stop() BackOff {
	return &backoff.StopBackOff{}
}

func Zero() BackOff {
	return &backoff.ZeroBackOff{}
}

func Const(d time.Duration) BackOff {
	if d > 0 {
		return backoff.NewConstantBackOff(d)
	} else {
		return &backoff.ZeroBackOff{}
	}
}

func TickerOf(b BackOff) *Ticker {
	return backoff.NewTicker(b)
}

func Retry(b BackOff, action func() error) error {
	err := backoff.Retry(action, b)
	return sderr.Wrap(err, "sdbackoff retry error")
}

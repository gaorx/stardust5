package sdconcur

import (
	"github.com/gaorx/stardust5/sdtime"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestDo(t *testing.T) {
	const (
		threads = 100
		n       = 200
		sleepMS = 10
	)
	var mtx sync.Mutex
	counter := 0
	start := sdtime.NowUnixMS()
	err := Do(0, lo.Range(threads), func(_, _ int) {
		for i := 0; i < n; i++ {
			sdtime.SleepMS(sleepMS)
			Lock(&mtx, func() {
				counter += 1
			})
		}
	})
	elapsed := sdtime.NowUnixMS() - start
	assert.NoError(t, err)
	assert.True(t, elapsed >= sleepMS*n && elapsed <= 2*sleepMS*n)
	assert.Equal(t, threads*n, counter)
}

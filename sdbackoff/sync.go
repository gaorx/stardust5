package sdbackoff

import (
	"sync"
	"time"
)

type syncBackOff struct {
	backOff BackOff
	mtx     sync.Mutex
}

func Synchronized(b BackOff) BackOff {
	return &syncBackOff{backOff: b}
}

func (b *syncBackOff) NextBackOff() time.Duration {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	return b.backOff.NextBackOff()
}

func (b *syncBackOff) Reset() {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.backOff.Reset()
}

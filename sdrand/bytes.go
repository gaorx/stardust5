package sdrand

import (
	"github.com/gaorx/stardust5/sdbytes"
	"math/rand"
)

func Bytes(n int) sdbytes.Slice {
	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Intn(16))
	}
	return b
}

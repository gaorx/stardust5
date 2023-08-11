package sdrand

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestSample(t *testing.T) {
	{
		n, total := 10, 1000000
		choices := lo.Range(n)
		c := newCounter[int]()
		for i := 0; i < total; i++ {
			v := Sample(choices...)
			c.inc(v)
		}
		assert.True(t, c.expectAll(1.0/float64(n), 0.01))
	}

	{
		n, total := 10, 1000000
		choices := lo.Range(n)
		c := newCounter[int]()
		for i := 0; i < total; i++ {
			some := Samples(choices, 3)
			for _, v := range some {
				c.inc(v)
			}
		}
		assert.True(t, c.expectAll(1.0/float64(n), 0.01))
	}

	{
		n, total := 10, 1000000
		choices := lo.Map(lo.Range(n), func(_ int, v int) W[int] {
			return W[int]{
				W: v + 1,
				V: v,
			}
		})
		c := newCounter[int]()
		for i := 0; i < total; i++ {
			v := SampleWeighted(choices...)
			c.inc(v)
		}
		for i := 0; i < n; i++ {
			assert.True(t, c.expectOne(i, (1.0/55.0)*float64(i+1), 0.01))
		}
	}
}

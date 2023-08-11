package sdrand

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBetween(t *testing.T) {
	{
		n, total := 10, 1000000
		c := newCounter[int]()
		for i := 0; i < total; i++ {
			v := IntBetween(111, 111+n)
			c.inc(v)
		}
		assert.True(t, c.expectAll(1.0/float64(n), 0.01))
	}

	{
		n, total := 10, 1000000
		c := newCounter[int64]()
		for i := 0; i < total; i++ {
			v := Int64Between(111, int64(111+n))
			c.inc(v)
		}
		assert.True(t, c.expectAll(1.0/float64(n), 0.01))
	}
}

type counter[T comparable] struct {
	data map[T]int64
}

func newCounter[T comparable]() counter[T] {
	return counter[T]{data: map[T]int64{}}
}

func (c counter[T]) size() int {
	return len(c.data)
}

func (c *counter[T]) inc(v T) {
	v1, ok := c.data[v]
	if ok {
		c.data[v] = v1 + 1
	} else {
		c.data[v] = 1
	}
}

func (c counter[T]) total() int64 {
	var total int64 = 0
	for _, counter := range c.data {
		total += counter
	}
	return total
}

func (c counter[T]) expectOne(v T, expected, epsilon float64) bool {
	total := c.total()
	n := c.data[v]
	actual := float64(n) / float64(total)
	return math.Abs(actual-expected) <= epsilon
}

func (c counter[T]) expectAll(expected, epsilon float64) bool {
	total := c.total()
	for _, counter := range c.data {
		actual := float64(counter) / float64(total)
		if math.Abs(actual-expected) > epsilon {
			return false
		}
	}
	return true
}

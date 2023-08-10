package sdparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	t0 := time.Now()
	for _, layout := range timeLayoutsForParse {
		s0 := t0.Format(layout)
		t1, err := Time(s0)
		assert.NoError(t, err)
		s1 := t1.Format(layout)
		assert.True(t, s0 == s1)
	}
}

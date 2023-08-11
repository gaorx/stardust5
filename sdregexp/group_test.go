package sdregexp

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestGroup(t *testing.T) {
	r := regexp.MustCompile(`^(?P<first>\w+)-(?P<last>\w+)$`)
	group := FindStringSubmatchGroup(r, "hello-world")
	assert.EqualValues(t, map[string]string{"first": "hello", "last": "world"}, group)
}

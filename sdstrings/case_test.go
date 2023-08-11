package sdstrings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCase(t *testing.T) {
	assert.Equal(t, "hello_world", ToSnakeL("HelloWorld"))
	assert.Equal(t, "HELLO_WORLD", ToSnakeU("HelloWorld"))
	assert.Equal(t, "HELLO_WORLD", ToSnakeU("hello_world"))

	assert.Equal(t, "hello-world", ToKebabL("HelloWorld"))
	assert.Equal(t, "HELLO-WORLD", ToKebabU("HelloWorld"))
	assert.Equal(t, "HELLO-WORLD", ToKebabU("hello_world"))
	assert.Equal(t, "HELLO-WORLD", ToKebabU("hello-world"))

	assert.Equal(t, "helloWorld", ToCamelL("hello_world"))
	assert.Equal(t, "HelloWorld", ToCamelU("hello-world"))
	assert.Equal(t, "HelloWorld", ToCamelU("helloWorld"))
}

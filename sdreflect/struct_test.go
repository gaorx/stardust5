package sdreflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructToMap(t *testing.T) {
	assert.Equal(
		t,
		map[string]any{"Name": "me", "Age": 44},
		StructToMap(s1{Name: "me", Age: 44}),
	)
}

func TestStructSelectFields(t *testing.T) {
	a := s1{
		Name: "me",
		Age:  44,
	}
	{
		var b s1
		err := StructSelectFields(&b, &a, []string{"Name"})
		assert.NoError(t, err)
		assert.Equal(t, "me", b.Name)
		assert.Equal(t, 0, b.Age)
	}
	{
		var b s1
		err := StructSelectFields(&b, &a, []string{"Age"})
		assert.NoError(t, err)
		assert.Equal(t, "", b.Name)
		assert.Equal(t, 44, b.Age)
	}
	{
		var b s1
		err := StructSelectFields(&b, &a, []string{"Age", "Name"})
		assert.NoError(t, err)
		assert.Equal(t, "me", b.Name)
		assert.Equal(t, 44, b.Age)
	}
}

type s1 struct {
	Name string
	Age  int
}

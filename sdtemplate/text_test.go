package sdtemplate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextTemplate(t *testing.T) {
	data := map[string]any{
		"Name": "<hello>",
	}
	_, err := Text.Exec(`{{.Name}`, data)
	assert.Error(t, err)
	s, err := Text.Exec(`{{.Name}}`, data)
	assert.NoError(t, err)
	assert.Equal(t, "<hello>", s)
}

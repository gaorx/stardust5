package sdtemplate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHtmlTemplate(t *testing.T) {
	data := map[string]any{
		"Name": "<hello>",
	}
	_, err := Html.Exec(`{{.Name}`, data)
	assert.Error(t, err)
	s, err := Html.Exec(`{{.Name}}`, data)
	assert.NoError(t, err)
	assert.Equal(t, "&lt;hello&gt;", s)
}

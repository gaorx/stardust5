package sdload

import (
	"github.com/gaorx/stardust5/sdfile"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	_ = sdfile.UseTempDir("", "", func(dirname string) {
		const text = "hello world"
		filename := filepath.Join(dirname, "a.txt")
		err := sdfile.WriteText(filename, text, 0600)
		assert.NoError(t, err)
		s, err := Text(filename)
		assert.NoError(t, err)
		assert.Equal(t, text, s)
		s, err = Text("file://" + filename)
		assert.NoError(t, err)
		assert.Equal(t, text, s)
	})
	s, err := Text("https://www.baidu.com")
	assert.NoError(t, err)
	assert.True(t, strings.Contains(s, "baidu"))
	s, err = Text("http://www.baidu.com")
	assert.NoError(t, err)
	assert.True(t, strings.Contains(s, "baidu"))
	_, err = Text("https://unavaiable_hostname/aaa.txt")
	assert.Error(t, err)
}

package sdfile

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRW(t *testing.T) {
	var tmpDir string
	err := UseTempDir("", "", func(dirname string) {
		tmpDir = dirname
		filename := filepath.Join(dirname, "a.txt")
		err1 := WriteText(filename, "hello", 0600)
		assert.NoError(t, err1)
		err1 = AppendText(filename, "world", 0600)
		assert.NoError(t, err1)
		s, err1 := ReadText(filename)
		assert.NoError(t, err1)
		assert.Equal(t, "helloworld", s)
	})
	assert.NoError(t, err)
	assert.NoDirExists(t, tmpDir)
}

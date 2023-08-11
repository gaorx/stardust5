package sdfile

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempDir(t *testing.T) {
	const text = "hello"
	var tmpDir string
	s, err := UseTempDirForResult("", "", func(dirname string) (string, error) {
		tmpDir = dirname
		filename := filepath.Join(dirname, "b.txt")
		err1 := WriteText(filename, text, 0600)
		assert.NoError(t, err1)
		s1, err1 := ReadText(filename)
		assert.NoError(t, err1)
		return s1, nil
	})
	assert.NoError(t, err)
	assert.NoDirExists(t, tmpDir)
	assert.Equal(t, text, s)
}

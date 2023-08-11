package sdexec

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	msg := "hello world!"
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	cmd, err := Parsef("echo '%s'", msg)
	assert.NoError(t, err)
	cmd.SetDir(homeDir)
	rr := cmd.RunResult()
	assert.NoError(t, rr.Err)
	assert.True(t, strings.Contains(rr.StdoutString(), msg))
}

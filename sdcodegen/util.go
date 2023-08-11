package sdcodegen

import (
	"os"
	"strings"

	"github.com/gaorx/stardust5/sderr"
)

func Expand(filename string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}
	filename = strings.Replace(filename, "~", homeDir, 1)
	return os.ExpandEnv(filename)
}

func TryExpand(filename string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", sderr.Wrap(err, "get home directory error")
	}
	filename = strings.Replace(filename, "~", homeDir, 1)
	os.Environ()
	var missing []string
	filename = os.Expand(filename, func(env string) string {
		v, ok := os.LookupEnv(env)
		if !ok {
			missing = append(missing, env)
		}
		return v
	})
	if len(missing) > 0 {
		return filename, sderr.Newf("no env $%s", missing[0])
	}
	return filename, nil
}

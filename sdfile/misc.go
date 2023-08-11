package sdfile

import (
	"os"
	"path/filepath"
)

// Directory

func IsDir(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

func BinDir() string {
	dir, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(dir)
}

func AbsByBin(filename string) string {
	if filename == "" {
		return BinDir()
	}
	if filepath.IsAbs(filename) {
		return filename
	}
	r, err := filepath.Abs(filepath.Join(BinDir(), filename))
	if err != nil {
		return ""
	}
	return r
}

// Exists

func Exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func FirstExists(filenames ...string) string {
	if len(filenames) == 0 {
		return ""
	}
	for _, filename := range filenames {
		if Exists(filename) {
			return filename
		}
	}
	return ""
}

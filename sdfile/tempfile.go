package sdfile

import (
	"github.com/samber/lo"
	"os"

	"github.com/gaorx/stardust5/sderr"
)

func UseTempFile(dir, pattern string, action func(*os.File)) error {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return sderr.WrapWith(err, "create temp file error", dir, pattern)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	action(f)
	return nil
}

func UseTempDir(dir, pattern string, action func(string)) error {
	name, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return sderr.WrapWith(err, "create temp dir error", dir, pattern)
	}
	defer func() {
		_ = os.RemoveAll(name)
	}()
	action(name)
	return nil
}

func UseTempFileForResult[R any](dir, pattern string, action func(*os.File) (R, error)) (R, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return lo.Empty[R](), sderr.WrapWith(err, "create temp file error", dir, pattern)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	if r, err := action(f); err != nil {
		return lo.Empty[R](), sderr.WrapWith(err, "call file action error", dir, pattern)
	} else {
		return r, nil
	}
}

func UseTempDirForResult[R any](dir, pattern string, action func(string) (R, error)) (R, error) {
	var empty R
	name, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return empty, sderr.WrapWith(err, "create temp dir error", dir, pattern)
	}
	defer func() {
		_ = os.RemoveAll(name)
	}()

	if r, err := action(name); err != nil {
		return empty, sderr.WrapWith(err, "call dir action error", dir, pattern)
	} else {
		return r, nil
	}
}

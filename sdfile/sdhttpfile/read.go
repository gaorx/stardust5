package sdhttpfile

import (
	"io"
	"net/http"

	"github.com/gaorx/stardust5/sderr"
)

func HttpReadBytes(hfs http.FileSystem, name string) ([]byte, error) {
	if hfs == nil {
		return nil, sderr.New("nil hfs")
	}
	f, err := hfs.Open(name)
	if err != nil {
		return nil, sderr.WrapWith(err, "open error", name)
	}
	defer func() { _ = f.Close() }()
	r, err := io.ReadAll(f)
	if err != nil {
		return nil, sderr.WrapWith(err, "read error", name)
	}
	return r, nil
}

func HttpReadText(hfs http.FileSystem, name string) (string, error) {
	b, err := HttpReadBytes(hfs, name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func HttpReadTextDef(hfs http.FileSystem, name, def string) string {
	s, err := HttpReadText(hfs, name)
	if err != nil {
		return def
	}
	return s
}

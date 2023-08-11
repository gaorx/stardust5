package sdfile

import (
	"io"
	"os"

	"github.com/gaorx/stardust5/sderr"
)

func ReadBytes(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, sderr.WrapWith(err, "read file error", filename)
	}
	return data, nil
}

func ReadBytesDef(filename string, def []byte) []byte {
	data, err := ReadBytes(filename)
	if err != nil {
		return def
	}
	return data
}

func WriteBytes(filename string, data []byte, perm os.FileMode) error {
	err := os.WriteFile(filename, data, perm)
	if err != nil {
		return sderr.WrapWith(err, "write file error", filename)
	}
	return nil
}

func AppendBytes(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return sderr.WrapWith(err, "open append error", filename)
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	_ = f.Close()
	return sderr.Wrap(err, "write for append error")
}

func ReadText(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", sderr.WrapWith(err, "read text error", filename)
	}
	return string(data), nil
}

func ReadTextDef(filename, def string) string {
	data, err := ReadText(filename)
	if err != nil {
		return def
	}
	return data
}

func WriteText(filename string, text string, perm os.FileMode) error {
	return WriteBytes(filename, []byte(text), perm)
}

func AppendText(filename string, text string, perm os.FileMode) error {
	return AppendBytes(filename, []byte(text), perm)
}

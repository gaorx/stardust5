package sdload

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"net/url"

	"github.com/BurntSushi/toml"
	"github.com/gaorx/stardust5/sderr"
)

func Bytes(loc string) ([]byte, error) {
	var scheme string
	u, err := url.Parse(loc)
	if err != nil {
		scheme = ""
	} else {
		scheme = u.Scheme
	}
	l, ok := loaders[scheme]
	if !ok {
		return nil, sderr.NewWith("unknown scheme", loc)
	}
	data, err := l.LoadBytes(loc)
	if err != nil {
		return nil, sderr.WrapWith(err, "load error", loc)
	}
	return data, nil
}

func Text(loc string) (string, error) {
	data, err := Bytes(loc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func JSON[T any](loc string) (T, error) {
	var empty, r T
	data, err := Bytes(loc)
	if err != nil {
		return empty, err
	}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return empty, sderr.Wrap(err, "parse json error")
	}
	return r, nil
}

func TOML[T any](loc string, v any) (T, error) {
	var empty, r T
	data, err := Bytes(loc)
	if err != nil {
		return empty, err
	}
	err = toml.Unmarshal(data, &r)
	if err != nil {
		return empty, sderr.Wrap(err, "parse TOML error")
	}
	return r, nil
}

func YAML[T any](loc string, v any) (T, error) {
	var empty, r T
	data, err := Bytes(loc)
	if err != nil {
		return empty, err
	}
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return empty, sderr.Wrap(err, "parse YAML error")
	}
	return r, nil
}

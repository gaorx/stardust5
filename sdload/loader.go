package sdload

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gaorx/stardust5/sderr"
)

// Loader

type Loader interface {
	LoadBytes(loc string) ([]byte, error)
}

// LoaderFunc

type LoaderFunc func(loc string) ([]byte, error)

func (f LoaderFunc) LoadBytes(loc string) ([]byte, error) {
	return f(loc)
}

// Loaders

var (
	loaders = map[string]Loader{
		"":      LoaderFunc(fileLoader),
		"file":  LoaderFunc(fileLoader),
		"http":  LoaderFunc(httpLoader),
		"https": LoaderFunc(httpLoader),
	}
)

func RegisterLoader(scheme string, loader Loader) {
	if scheme == "" {
		panic(sderr.New("no scheme"))
	}
	if loader == nil {
		panic(sderr.New("nil loader"))
	}
	loaders[scheme] = loader
}

// default loader

func fileLoader(loc string) ([]byte, error) {
	loc = strings.TrimPrefix(loc, "file://")
	data, err := os.ReadFile(loc)
	if err != nil {
		return nil, sderr.Wrap(err, "read file error")
	}
	return data, nil
}

func httpLoader(loc string) ([]byte, error) {
	resp, err := (&http.Client{Timeout: 7 * time.Second}).Get(loc)
	if err != nil {
		return nil, sderr.Wrap(err, "http get error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, sderr.NewWith("response HTTP status error", resp.StatusCode, loc)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, sderr.Wrap(err, "read http response body error")
	}
	return data, nil
}

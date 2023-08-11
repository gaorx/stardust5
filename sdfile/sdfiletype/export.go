package sdfiletype

import (
	"io"

	"github.com/gaorx/stardust5/sderr"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type (
	Type = types.Type
)

// data

func Match(data []byte) (Type, error) {
	t, err := filetype.Match(data)
	if err != nil {
		return Type{}, sderr.Wrap(err, "match error")
	}
	return t, nil
}

func MatchMime(data []byte, def string) string {
	t, err := Match(data)
	if err != nil {
		return def
	}
	return t.MIME.Value
}

func MatchExt(data []byte, def string) string {
	t, err := Match(data)
	if err != nil {
		return def
	}
	return t.Extension
}

// reader

func MatchReader(r io.Reader) (Type, error) {
	t, err := filetype.MatchReader(r)
	if err != nil {
		return Type{}, sderr.Wrap(err, "match reader error")
	}
	return t, nil
}

func MatchReaderMime(r io.Reader, def string) string {
	t, err := MatchReader(r)
	if err != nil {
		return def
	}
	return t.MIME.Value
}

func MatchReaderExt(r io.Reader, def string) string {
	t, err := MatchReader(r)
	if err != nil {
		return def
	}
	return t.Extension
}

// file

func MatchFile(filename string) (Type, error) {
	t, err := filetype.MatchFile(filename)
	if err != nil {
		return Type{}, sderr.Wrap(err, "match file error")
	}
	return t, nil
}

func MatchFileMime(filename, def string) string {
	t, err := MatchFile(filename)
	if err != nil {
		return def
	}
	return t.MIME.Value
}

func MatchFileExt(filename, def string) string {
	t, err := MatchFile(filename)
	if err != nil {
		return def
	}
	return t.Extension
}

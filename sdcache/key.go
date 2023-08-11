package sdcache

import (
	"github.com/gaorx/stardust5/sderr"
	"strconv"
	"strings"
)

type Key interface {
	EncodeKey(k any) (string, error)
	DecodeKey(s string) (any, error)
	PrefixForClear() string
}

type FuncKey struct {
	Encode func(k any) (string, error)
	Decode func(s string) (any, error)
	Prefix string
}

func (fk FuncKey) EncodeKey(k any) (string, error) {
	if k == nil {
		return "", sderr.New("nil key")
	}
	s, err := fk.Encode(k)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return s, nil
}

func (fk FuncKey) DecodeKey(s string) (any, error) {
	k, err := fk.Decode(s)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return k, nil
}

func (fk FuncKey) PrefixForClear() string {
	return fk.Prefix
}

type StrKey struct{ Prefix string }
type Int64Key struct{ Prefix string }

func (sk StrKey) EncodeKey(k any) (string, error) {
	if k == nil {
		return "", sderr.New("nil key")
	}
	k1, ok := k.(string)
	if !ok {
		return "", sderr.New("key type error")
	}
	return sk.Prefix + k1, nil
}

func (sk StrKey) DecodeKey(s string) (any, error) {
	prefix := sk.Prefix
	if prefix != "" {
		if !strings.HasPrefix(s, prefix) {
			return nil, sderr.New("key prefix error")
		}
		return s[len(prefix):], nil
	} else {
		return s, nil
	}
}

func (sk StrKey) PrefixForClear() string {
	return sk.Prefix
}

func (ik Int64Key) EncodeKey(k any) (string, error) {
	if k == nil {
		return "", sderr.New("nil key")
	}
	k1, ok := k.(int64)
	if !ok {
		return "", sderr.New("key type error")
	}
	return ik.Prefix + strconv.FormatInt(k1, 10), nil
}

func (ik Int64Key) DecodeKey(s string) (any, error) {
	prefix := ik.Prefix
	if prefix != "" {
		if !strings.HasPrefix(s, prefix) {
			return nil, sderr.New("key prefix error")
		}
		s = s[len(prefix):]
	}
	k1, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, sderr.New("parse int64 key error")
	}
	return k1, nil
}

func (ik Int64Key) PrefixForClear() string {
	return ik.Prefix
}

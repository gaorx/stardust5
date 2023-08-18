package sdobjectstore

import (
	"github.com/gaorx/stardust5/sderr"
)

var (
	Discard Interface = discard{}
)

type discard struct {
}

func (_ discard) Store(src Source, objectName string) (*Target, error) {
	if src == nil {
		return nil, sderr.New("nil source")
	}
	return &Target{
		Typ: DiscardTarget,
	}, nil
}

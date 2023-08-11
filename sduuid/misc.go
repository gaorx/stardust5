package sduuid

import (
	"github.com/gaorx/stardust5/sdbytes"
	"github.com/gofrs/uuid"
)

type (
	UUID = uuid.UUID
)

func Encode(id UUID) sdbytes.Slice {
	return id.Bytes()
}

func NewV1() sdbytes.Slice {
	v, err := uuid.NewV1()
	if err != nil {
		return nil
	}
	return v.Bytes()
}

func NewV4() sdbytes.Slice {
	v, err := uuid.NewV4()
	if err != nil {
		return nil
	}
	return v.Bytes()
}

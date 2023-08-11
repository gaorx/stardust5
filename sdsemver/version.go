package sdsemver

import (
	"fmt"

	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdparse"
	"github.com/gaorx/stardust5/sdstrings"
)

type V struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

var (
	ErrParse = sderr.New("sdsemver parse version error")
)

const numLimit = 10000

func New(major, minor, patch int) V {
	return V{Major: major, Minor: minor, Patch: patch}
}

func (v V) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v V) IsEmpty() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

func (v V) IsValidate() bool {
	return (v.Major >= 0 && v.Major < numLimit) &&
		(v.Minor >= 0 && v.Minor < numLimit) &&
		(v.Patch >= 0 && v.Patch < numLimit)
}

func (v V) ToInt() int64 {
	return int64((v.Major * numLimit * numLimit) + (v.Minor * numLimit) + v.Patch)
}

func (v V) Equal(major, minor, patch int) bool {
	return v.Major == major && v.Minor == minor && v.Patch == patch
}

func Parse(s string) (V, error) {
	if s == "" {
		return V{}, sderr.Wrap(ErrParse, "sdsemver parse empty")
	}
	majorStr, minorStr, patchStr := sdstrings.Split3s(s, ".")
	major, err := sdparse.Int(sdstrings.EmptyAs(majorStr, "0"))
	if err != nil {
		return V{}, sderr.Wrap(ErrParse, "sdsemver parse major error")
	}
	minor, err := sdparse.Int(sdstrings.EmptyAs(minorStr, "0"))
	if err != nil {
		return V{}, sderr.Wrap(ErrParse, "sdsemver parse minor error")
	}
	patch, err := sdparse.Int(sdstrings.EmptyAs(patchStr, "0"))
	if err != nil {
		return V{}, sderr.Wrap(ErrParse, "sdsemver parse patch error")
	}
	if !(major >= 0 && major < numLimit) {
		return V{}, sderr.Wrap(ErrParse, "sdsemver illegal major")
	}
	if !(minor >= 0 && minor < numLimit) {
		return V{}, sderr.Wrap(ErrParse, "sdsemver illegal minor")
	}
	if !(patch >= 0 && patch < numLimit) {
		return V{}, sderr.Wrap(ErrParse, "sdsemver illegal patch")
	}
	return V{Major: major, Minor: minor, Patch: patch}, nil
}

func FromInt(i int64) (V, error) {
	major := i / (numLimit * numLimit)
	minor := (i - major*numLimit*numLimit) / numLimit
	patch := i - major*numLimit*numLimit - minor*numLimit
	if !(major >= 0 && major < numLimit) {
		return V{}, sderr.Wrap(ErrParse, "sdsemver illegal major")
	}
	if !(minor >= 0 && minor < numLimit) {
		return V{}, sderr.Wrap(ErrParse, "sdsemver illegal minor")
	}
	if !(patch >= 0 && patch < numLimit) {
		return V{}, sderr.Wrap(ErrParse, "sdsemver illegal patch")
	}
	return V{Major: int(major), Minor: int(minor), Patch: int(patch)}, nil
}

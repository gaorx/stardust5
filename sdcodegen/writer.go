package sdcodegen

import (
	"io/fs"
)

type Writer interface {
	Len() int
	IsEmpty() bool
	Filename() string
	SetFilename(filename string) Writer
	Perm() fs.FileMode
	SetPerm(perm fs.FileMode) Writer
	Overwrite() bool
	SetOverwrite(b bool) Writer
	Modify(func(s string) string) Writer
	WritePlaceholder(p *Placeholder) Writer
	UsePlaceholder(name string, f func(p *Placeholder)) Writer
	P(s string) Writer
	F(format string, a ...any) Writer
	T(template string, data any) Writer
	L(s string) Writer
	FL(format string, a ...any) Writer
	TL(template string, data any) Writer
	NL() Writer
	I(n int) Writer
	Repeat(s string, n int) Writer
	If(cond bool, s string) Writer
	IfFunc(cond bool, f func()) Writer
	IfElse(cond bool, s, els string) Writer
	IfElseFunc(cond bool, f, els func()) Writer
	Join(l []string, sep, lastSep string) Writer
	JoinPairs(pairs []Pair, space, sep, lastSep string) Writer
}

type Pair struct {
	F, S string
}

type Placeholder struct {
	Name        string
	Data        any
	Expand      func(w Writer, data any)
	placeholder string
}

package sdcodegen

import (
	"fmt"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdfile"
	"github.com/gaorx/stardust5/sdrand"
	"github.com/gaorx/stardust5/sdtemplate"
	"github.com/samber/lo"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Buffer struct {
	options      BufferOptions
	data         strings.Builder
	filename     string
	perm         fs.FileMode
	overwrite    bool
	placeholders []*Placeholder
}

type BufferOptions struct {
	Newline string
	Indent  string
}

type Action string

const (
	I = Action("ignore")
	W = Action("write")
	C = Action("create")
)

func NewBuffer(opts BufferOptions) *Buffer {
	if opts.Newline == "" {
		opts.Newline = "\n"
	}
	if opts.Indent == "" {
		opts.Indent = "  "
	}
	return &Buffer{
		options:   opts,
		overwrite: true,
	}
}

func (b *Buffer) Clear() {
	b.data.Reset()
}

func (b *Buffer) String() string {
	if b == nil {
		return ""
	}
	return b.clone().expandPlaceholders().data.String()
}

func (b *Buffer) Bytes() []byte {
	if b == nil {
		return nil
	}
	return []byte(b.clone().expandPlaceholders().data.String())
}

func (b *Buffer) Save(root string, mkdirs bool, logger func(action Action, fn, absFn string)) error {
	var err error
	var absFn string

	if root == "" || root == "." || root == "."+string(filepath.Separator) {
		root, err = os.Getwd()
		if err != nil {
			return sderr.Wrap(err, "save buffers error")
		}
	}
	root = strings.TrimSuffix(root, string(filepath.Separator))

	fn := b.filename
	fn, err = TryExpand(fn)
	if err != nil {
		return sderr.Wrap(err, "expand buffer filename error")
	}

	if filepath.IsAbs(fn) {
		absFn = fn
	} else {
		root, err = TryExpand(root)
		if err != nil {
			return sderr.Wrap(err, "expand root error")
		}
		absFn, err = filepath.Abs(filepath.Join(root, fn))
		if err != nil {
			return sderr.WrapWith(err, "get absolute filename error", fn)
		}
	}

	var action Action
	if sdfile.Exists(absFn) {
		if b.overwrite {
			action = W
		} else {
			action = I
		}
	} else {
		action = C
	}
	if action == W || action == C {
		if mkdirs {
			dir := filepath.Dir(absFn)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return sderr.WrapWith(err, "mkdirs error", dir)
			}
		}
		perm := b.perm
		if perm == 0 {
			perm = 0644
		}
		err = os.WriteFile(absFn, b.clone().expandPlaceholders().Bytes(), perm)
		if err != nil {
			return sderr.WrapWith(err, "save file error", absFn)
		}
	}
	if logger != nil {
		logger(action, b.filename, absFn)
	}
	return nil
}

func (b *Buffer) expandPlaceholders() *Buffer {
	if len(b.placeholders) <= 0 {
		return b
	}
	data := b.data.String()
	for _, p0 := range b.placeholders {
		expanded := ""
		if p0.Expand != nil {
			expandedBuff := NewBuffer(b.options)
			p0.Expand(expandedBuff, p0.Data)
			expanded = expandedBuff.String()
		} else {
			if p0.Data != nil {
				if s, ok := p0.Data.(string); ok {
					expanded = s
				} else {
					expanded = fmt.Sprintf("%v", p0.Data)
				}
			}
		}
		data = strings.Replace(data, p0.placeholder, expanded, -1)
	}
	b.data.Reset()
	b.data.WriteString(data)
	return b
}

// functions

func (b *Buffer) Len() int {
	return b.data.Len()
}

func (b *Buffer) IsEmpty() bool {
	return b.Len() <= 0
}

func (b *Buffer) Filename() string {
	return b.filename
}

func (b *Buffer) SetFilename(filename string) Writer {
	b.filename = filename
	return b
}

func (b *Buffer) Perm() fs.FileMode {
	return b.perm
}

func (b *Buffer) SetPerm(perm fs.FileMode) Writer {
	b.perm = perm
	return b
}

func (b *Buffer) Overwrite() bool {
	return b.overwrite
}

func (b *Buffer) SetOverwrite(overwrite bool) Writer {
	b.overwrite = overwrite
	return b
}

func (b *Buffer) Modify(f func(s string) string) Writer {
	if f == nil {
		return b
	}
	data := b.data.String()
	data1 := f(data)
	b.data.Reset()
	b.data.WriteString(data1)
	return b
}

func (b *Buffer) WritePlaceholder(p *Placeholder) Writer {
	if p == nil {
		panic(sderr.New("nil placeholder"))
	}
	if p.Name == "" {
		panic(sderr.New("no placeholder name"))
	}
	for _, p0 := range b.placeholders {
		if p0.Name == p.Name {
			panic(sderr.NewWith("duplicated placeholder", p0.Name))
		}
	}
	p1 := *p
	p1.placeholder = fmt.Sprintf("<<<<%s:%s>>>>", p.Name, sdrand.String(8, sdrand.LowerCaseAlphanumericCharset))
	b.placeholders = append(b.placeholders, &p1)
	b.data.WriteString(p1.placeholder)
	return b
}

func (b *Buffer) UsePlaceholder(name string, f func(p *Placeholder)) Writer {
	if f == nil {
		return b
	}

	for _, p0 := range b.placeholders {
		if p0.Name == name {
			f(p0)
		}
	}
	return b
}

func (b *Buffer) P(s string) Writer {
	b.data.WriteString(s)
	return b
}

func (b *Buffer) F(format string, a ...any) Writer {
	return b.P(fmt.Sprintf(format, a...))
}

func (b *Buffer) T(template string, data any) Writer {
	s := lo.Must(sdtemplate.Text.Exec(template, data))
	return b.P(s)
}

func (b *Buffer) L(s string) Writer {
	b.P(s)
	return b.P(b.options.Newline)
}

func (b *Buffer) FL(format string, a ...any) Writer {
	return b.F(format, a...).NL()
}

func (b *Buffer) TL(template string, data any) Writer {
	return b.T(template, data).NL()
}

func (b *Buffer) NL() Writer {
	return b.P(b.options.Newline)
}

func (b *Buffer) Repeat(s string, n int) Writer {
	if n < 0 {
		n = 0
	}
	return b.P(strings.Repeat(s, n))
}

func (b *Buffer) I(n int) Writer {
	if n < 0 {
		n = 0
	}
	return b.P(strings.Repeat(b.options.Indent, n))
}

func (b *Buffer) If(cond bool, s string) Writer {
	if cond {
		b.P(s)
	}
	return b
}

func (b *Buffer) IfFunc(cond bool, f func()) Writer {
	if cond && f != nil {
		f()
	}
	return b
}

func (b *Buffer) IfElse(cond bool, s, els string) Writer {
	if cond {
		b.P(s)
	} else {
		b.P(els)
	}
	return b
}

func (b *Buffer) IfElseFunc(cond bool, f, els func()) Writer {
	if cond {
		if f != nil {
			f()
		}
	} else {
		if els != nil {
			els()
		}
	}
	return b
}

func (b *Buffer) Join(l []string, sep, lastSep string) Writer {
	n := len(l)
	for i, elem := range l {
		b.P(elem)
		if i < n-1 {
			b.P(sep)
		} else {
			b.P(lastSep)
		}
	}
	return b
}

func (b *Buffer) JoinPairs(pairs []Pair, space, sep, lastSep string) Writer {
	n := len(pairs)
	for i, pair := range pairs {
		b.P(pair.F).P(space).P(pair.S)
		if i < n-1 {
			b.P(sep)
		} else {
			b.P(lastSep)
		}
	}
	return b
}

func (b *Buffer) clone() *Buffer {
	var builder1 strings.Builder
	builder1.WriteString(b.data.String())
	var placeholders1 []*Placeholder
	for _, p0 := range b.placeholders {
		placeholders1 = append(placeholders1, &Placeholder{
			Name:        p0.Name,
			Data:        p0.Data,
			Expand:      p0.Expand,
			placeholder: p0.placeholder,
		})
	}
	return &Buffer{
		options:      b.options,
		data:         builder1,
		filename:     b.filename,
		perm:         b.perm,
		overwrite:    b.overwrite,
		placeholders: placeholders1,
	}
}

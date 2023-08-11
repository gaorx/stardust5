package sdcodegen

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/gaorx/stardust5/sderr"
)

type Buffers struct {
	buffers []*Buffer
}

func NewBuffers() *Buffers {
	return &Buffers{}
}

func (bs *Buffers) String() string {
	return strings.Join(bs.Filenames(), "\n")
}

func (bs *Buffers) Dump(out io.Writer) {
	p := func(ss ...string) {
		for _, s := range ss {
			_, _ = out.Write([]byte(s))
		}
	}
	for _, b := range bs.buffers {
		p(b.filename, "\n")
		p(strings.Repeat("-", 60), "\n")
		lines := strings.Split(b.String(), "\n")
		for _, line := range lines {
			p("> ", strings.TrimSuffix(line, "\r"), "\n")
		}
		p("\n")
	}
}

func (bs *Buffers) Has(filename string) bool {
	return bs.Get(filename) != nil
}

func (bs *Buffers) Get(filename string) *Buffer {
	for _, b := range bs.buffers {
		if b.filename == filename {
			return b
		}
	}
	return nil
}

func (bs *Buffers) Data(filename string) string {
	buff := bs.Get(filename)
	if buff == nil {
		return ""
	}
	return buff.String()
}

func (bs *Buffers) Find(pattern string) []*Buffer {
	var r []*Buffer
	for _, b := range bs.buffers {
		ok, err := filepath.Match(pattern, b.filename)
		if err == nil && ok {
			r = append(r, b)
		}
	}
	return r
}

func (bs *Buffers) Filenames() []string {
	var filenames []string
	for _, b := range bs.buffers {
		filenames = append(filenames, b.filename)
	}
	return filenames
}

func (bs *Buffers) Save(root string, logger func(action Action, fn, absFn string)) error {
	for _, b := range bs.buffers {
		err := b.Save(root, true, logger)
		if err != nil {
			return sderr.WithStack(err)
		}
	}
	return nil
}

func (bs *Buffers) Put(filename string, f func(w Writer)) *Buffer {
	buff := bs.Open(filename)
	if f != nil {
		f(buff)
	}
	return buff
}

func (bs *Buffers) Add(filename string, f func(w Writer)) *Buffer {
	buff := bs.Append(filename)
	if f != nil {
		f(buff)
	}
	return buff
}

func (bs *Buffers) Open(filename string) *Buffer {
	buff, err := bs.open(filename)
	if err != nil {
		panic(err)
	}
	buff.Clear()
	return buff
}

func (bs *Buffers) Append(filename string) *Buffer {
	buff, err := bs.open(filename)
	if err != nil {
		panic(err)
	}
	return buff
}

func (bs *Buffers) open(filename string) (*Buffer, error) {
	if filename == "" {
		return nil, sderr.New("no filename")
	}

	buff := bs.Get(filename)
	if buff != nil {
		return buff, nil
	}

	buff = NewBuffer(getBufferOptions(filename))
	buff.SetFilename(filename)
	bs.buffers = append(bs.buffers, buff)
	return buff, nil
}

func getBufferOptions(filename string) BufferOptions {
	const newline = "\n"
	base := filepath.Base(filename)
	if base == "Makefile" || base == "Makefile.in" {
		return BufferOptions{Newline: newline, Indent: "\t"}
	}
	ext := strings.ToLower(filepath.Ext(base))
	switch ext {
	case ".go":
		return BufferOptions{Newline: newline, Indent: "\t"}
	case ".htm", ".html", ".json", ".js", ".jsx", ".ts", ".tsx", ".coffee", ".vue", ".astro", ".css", ".scss", ".sass", ".less":
		return BufferOptions{Newline: newline, Indent: "  "}
	case ".yaml", ".yml", ".toml", ".ini", ".properties":
		return BufferOptions{Newline: newline, Indent: "  "}
	case ".md":
		return BufferOptions{Newline: newline, Indent: "  "}
	case ".py":
		return BufferOptions{Newline: newline, Indent: "    "}
	case ".c", ".h", ".cpp", ".cxx", ".hpp", ".hxx", ".zig":
		return BufferOptions{Newline: newline, Indent: "    "}
	case ".java":
		return BufferOptions{Newline: newline, Indent: "    "}
	case ".kt":
		return BufferOptions{Newline: newline, Indent: "  "}
	case ".rs":
		return BufferOptions{Newline: newline, Indent: "  "}
	default:
		return BufferOptions{Newline: newline, Indent: "  "}
	}
}

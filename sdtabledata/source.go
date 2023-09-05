package sdtabledata

import (
	"io/fs"
	"os"
	"path"
)

type Source struct {
	Root    fs.FS
	Dir     string
	trimmed bool
}

func (src Source) IsNil() bool {
	return src.Root == nil
}

func (src Source) Trim() Source {
	if src.trimmed {
		return src
	}
	return Source{
		Root:    src.Root,
		Dir:     trimDir(src.Dir),
		trimmed: true,
	}
}

func (src Source) IsTrimmed() bool {
	return src.trimmed
}

func (src Source) Sub(subs ...string) Source {
	src1 := src.Trim()
	subs1 := []string{src1.Dir}
	for _, sub := range subs {
		if sub != "" && sub != "." {
			subs1 = append(subs1, sub)
		}
	}
	return Source{
		Root: src1.Root,
		Dir:  path.Join(subs1...),
	}.Trim()
}

func Dir(dirname string, subs ...string) Source {
	return Source{
		Root: os.DirFS(dirname),
		Dir:  path.Join(subs...),
	}
}

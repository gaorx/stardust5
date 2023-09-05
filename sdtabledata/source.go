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

func (src Source) Sub(sub string) Source {
	if sub == "" || sub == "." {
		return src
	}
	src1 := src.Trim()
	return Source{
		Root: src1.Root,
		Dir:  path.Join(src1.Dir, sub),
	}.Trim()
}

func Dir(dirname string, subs ...string) Source {
	return Source{
		Root: os.DirFS(dirname),
		Dir:  path.Join(subs...),
	}
}

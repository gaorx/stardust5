package sdtabledata

import (
	"io/fs"
	"os"
	"path"
)

type Source struct {
	Root fs.FS
	Sub  string
}

func (src Source) IsNil() bool {
	return src.Root == nil
}

func (src Source) Trim() Source {
	return Source{
		Root: src.Root,
		Sub:  trimDir(src.Sub),
	}
}

func Dir(dirname string, subs ...string) Source {
	return Source{
		Root: os.DirFS(dirname),
		Sub:  path.Join(subs...),
	}
}

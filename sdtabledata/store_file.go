package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdobjectstore"
	"github.com/samber/lo"
	"io/fs"
)

type StoreFile = func(fsys fs.FS, fn string) (string, error)

// object store

type ObjectStoreOptions struct {
	ObjectName string
	HttpUrl    bool
}

func ObjectStore(store sdobjectstore.Store, opts *ObjectStoreOptions) StoreFile {
	opts1 := lo.FromPtr(opts)
	return func(fsys fs.FS, fn string) (string, error) {
		if store.IsNil() {
			return "", nil
		}
		target, err := store.StoreFileFS(fsys, fn, opts1.ObjectName)
		if err != nil {
			return "", sderr.WrapWith(err, "store column file error", fn)
		}
		if opts1.HttpUrl {
			return target.Url(), nil
		} else {
			return target.HttpsUrl(), nil
		}
	}
}

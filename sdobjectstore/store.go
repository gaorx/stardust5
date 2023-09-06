package sdobjectstore

import (
	"github.com/gaorx/stardust5/sderr"
	"io/fs"
)

type Interface interface {
	Store(src Source, objectName string) (*Target, error)
}

type Store struct {
	Interface
}

func (s Store) IsNil() bool {
	return s.Interface == nil
}

func (s Store) StoreFile(filename, objectName string) (*Target, error) {
	return s.Store(File(filename, ""), objectName)
}

func (s Store) StoreData(data []byte, objectName string) (*Target, error) {
	return s.Store(Bytes(data, ""), objectName)
}

func (s Store) StoreFileFS(fsys fs.FS, fn string, objectName string) (*Target, error) {
	data, err := fs.ReadFile(fsys, fn)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return s.Store(Bytes(data, ""), objectName)
}

func Dir(root string) Store {
	return Store{dir{root}}
}

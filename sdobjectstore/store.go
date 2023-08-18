package sdobjectstore

type Interface interface {
	Store(src Source, objectName string) (*Target, error)
}

type Store struct {
	Interface
}

func (s Store) StoreFile(filename, objectName string) (*Target, error) {
	return s.Store(File(filename, ""), objectName)
}

func (s Store) StoreData(data []byte, objectName string) (*Target, error) {
	return s.Store(Bytes(data, ""), objectName)
}

func Dir(root string) Store {
	return Store{dir{root}}
}

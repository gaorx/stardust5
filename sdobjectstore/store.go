package sdobjectstore

type Store interface {
	Store(src Source, objectName string) (*Target, error)
}

package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
)

func LoadAllT[T any](src Source, mapper func(sdjson.Object) (T, error), opts *LoadOptions) ([]T, error) {
	rows1, err := LoadAll(src, opts)
	if err != nil {
		return nil, err
	}
	return mapRows(rows1, mapper)
}

func LoadSomeT[T any](src Source, rows []string, mapper func(sdjson.Object) (T, error), opts *LoadOptions) ([]T, error) {
	rows1, err := LoadSome(src, rows, opts)
	if err != nil {
		return nil, err
	}
	return mapRows(rows1, mapper)
}

func LoadOneT[T any](src Source, dir string, row string, mapper func(sdjson.Object) (T, error), opts *LoadOptions) (T, error) {
	row1, err := LoadOne(src, row, opts)
	if err != nil {
		var empty T
		return empty, err
	}
	if mapper == nil {
		mapper = defaultMapper[T]()
	}
	return mapper(row1)
}

func mapRows[T any](rows1 []sdjson.Object, mapper func(sdjson.Object) (T, error)) ([]T, error) {
	if mapper == nil {
		mapper = defaultMapper[T]()
	}
	rows2 := make([]T, 0, len(rows1))
	for _, row1 := range rows1 {
		if row2, err := mapper(row1); err != nil {
			return nil, sderr.WithStack(err)
		} else {
			rows2 = append(rows2, row2)
		}
	}
	return rows2, nil
}

func defaultMapper[T any]() func(sdjson.Object) (T, error) {
	return func(row sdjson.Object) (T, error) {
		return sdjson.ObjectToStruct[T](row)
	}
}

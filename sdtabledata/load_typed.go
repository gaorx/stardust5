package sdtabledata

func LoadAllT[T any](src Source, opts *LoadOptions) ([]T, error) {
	var rows1 []T
	err := LoadAll(src, &rows1, opts)
	if err != nil {
		return nil, err
	}
	return rows1, nil
}

func LoadSomeT[T any](src Source, rows []string, opts *LoadOptions) ([]T, error) {
	var rows1 []T
	err := LoadSome(src, &rows1, rows, opts)
	if err != nil {
		return nil, err
	}
	return rows1, nil
}

func LoadOneT[T any](src Source, dir string, row string, opts *LoadOptions) (T, error) {
	var row1 T
	err := LoadOne(src, row, &row1, opts)
	if err != nil {
		var empty T
		return empty, err
	}
	return row1, nil
}

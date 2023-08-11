package sdtabledata

func LoadAllT[T any](dir string, opts LoadOptions) ([]T, error) {
	var rows1 []T
	err := LoadAll(dir, &rows1, opts)
	if err != nil {
		return nil, err
	}
	return rows1, nil
}

func LoadSomeT[T any](dir string, rows []string, opts LoadOptions) ([]T, error) {
	var rows1 []T
	err := LoadSome(dir, &rows1, rows, opts)
	if err != nil {
		return nil, err
	}
	return rows1, nil
}

func LoadOneT[T any](dir string, row string, opts LoadOptions) (T, error) {
	var row1 T
	err := LoadOne(dir, row, &row1, opts)
	if err != nil {
		var empty T
		return empty, err
	}
	return row1, nil
}

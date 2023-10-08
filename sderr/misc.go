package sderr

func Combine(errs []error) error {
	var merged error = nil
	for _, err := range errs {
		if err != nil {
			merged = Append(merged, err)
		}
	}
	return merged
}

func Scatter(err error) []error {
	if err == nil {
		return nil
	}
	if err1, ok := AsT[*MultipleError](err); ok && err1 != nil {
		return err1.Errors
	}
	return []error{err}
}

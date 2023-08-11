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

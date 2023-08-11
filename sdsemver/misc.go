package sdsemver

func ToInt(s string) (int64, error) {
	v, err := Parse(s)
	if err != nil {
		return 0, err
	}
	return v.ToInt(), nil
}

func ToString(i int64) (string, error) {
	v, err := FromInt(i)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

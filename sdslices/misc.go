package sdslices

func Ensure[T any](s []T) []T {
	if s == nil {
		return make([]T, 0)
	}
	return s
}

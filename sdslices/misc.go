package sdslices

func Ensure[S ~[]T, T any](s S) S {
	if s == nil {
		return make(S, 0)
	}
	return s
}

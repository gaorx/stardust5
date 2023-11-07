package sdcheck

import (
	"slices"
)

func In[T comparable, C ~[]T](v T, available C, message any) CheckerFunc {
	return func() error {
		if !slices.Contains[C, T](available, v) {
			return errorOf(message)
		}
		return nil
	}
}

func NotIn[T comparable, C ~[]T](v T, available C, message any) CheckerFunc {
	return func() error {
		if slices.Contains[C, T](available, v) {
			return errorOf(message)
		}
		return nil
	}
}

func HasKey[K comparable, V any, M ~map[K]V](k K, m M, message any) CheckerFunc {
	return func() error {
		if _, ok := m[k]; !ok {
			return errorOf(message)
		}
		return nil
	}
}

func NotHasKey[K comparable, V any, M ~map[K]V](k K, m M, message any) CheckerFunc {
	return func() error {
		if _, ok := m[k]; ok {
			return errorOf(message)
		}
		return nil
	}
}

package util

func Merge[T any](a, b []T) []T {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	result := make([]T, len(a)+len(b))
	n := copy(result, a)
	copy(result[n:], b)
	return result
}

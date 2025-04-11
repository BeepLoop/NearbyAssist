package utils

func Ternary[T any](condition bool, good T, fallback T) T {
	if condition {
		return good
	}

	return fallback
}

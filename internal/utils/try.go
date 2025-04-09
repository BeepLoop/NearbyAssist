package utils

func Try[T any](value T, err error) T {
	if err != nil {
		return *new(T)
	}

	return value
}

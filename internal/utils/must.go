package utils

func Must[T any](value T, err error) T {
	if err != nil {
		panic("Error in Must command: " + err.Error())
	}

	return value
}

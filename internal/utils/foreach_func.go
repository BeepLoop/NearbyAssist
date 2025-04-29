package utils

func ForEach[T any](items []T, fn func(T)) {
	for _, item := range items {
		fn(item)
	}
}

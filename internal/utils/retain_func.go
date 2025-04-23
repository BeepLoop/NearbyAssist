package utils

import "iter"

func Retain[T any](input []T, fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range input {
			if !fn(item) {
				continue
			}

			if !yield(item) {
				return
			}
		}
	}
}

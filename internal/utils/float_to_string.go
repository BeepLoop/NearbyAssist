package utils

import "strconv"

func Float64ToString(input float64) string {
	return strconv.FormatFloat(input, 'f', -1, 64)
}

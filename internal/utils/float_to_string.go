package utils

import "strconv"

func FloatToString(input float64) string {
	return strconv.FormatFloat(input, 'f', -1, 64)
}

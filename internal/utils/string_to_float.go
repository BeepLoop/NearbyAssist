package utils

import "strconv"

func StringToFloatElseZero(input string) float64 {
	v, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0.0
	}

	return v
}

package utils

import "strconv"

func StringToFloat64ElseZero(input string) float64 {
	v, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0.0
	}

	return v
}

func StringToFloat32ElseZero(input string) float32 {
	v, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0.0
	}

	return float32(v)
}

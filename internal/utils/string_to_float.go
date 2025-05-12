package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func StringToFloat64ElseZero(input string) float64 {
	v, err := strconv.ParseFloat(strings.ReplaceAll(input, ",", ""), 64)
	if err != nil {
		fmt.Println("error convert: ", err.Error())
		return 0.0
	}

	return v
}

func StringToFloat32ElseZero(input string) float32 {
	return float32(StringToFloat64ElseZero(input))
}

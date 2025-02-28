package utils

import (
	"strconv"
)

func PercentageIncrease(current, previous int) float64 {
	increase := float64(current) - float64(previous)
	if increase == 0 {
		return 0
	}

	return (increase / float64(previous)) * 100
}

func PercentageDecrease(current, previous int) float64 {
	decrease := float64(previous) - float64(current)
	if decrease == 0 {
		return 0
	}

	return (decrease / float64(previous)) * 100
}

func ToPercentageString(amount float64) string {
	return strconv.FormatFloat(amount, 'f', -1, 64) + "%"
}

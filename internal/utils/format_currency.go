package utils

import (
	"fmt"
	"strconv"
)

func FormatCurrency(input float64) string {
	amount := strconv.FormatFloat(input, 'f', -1, 64)

	return fmt.Sprintf("₱ %s", amount)
}

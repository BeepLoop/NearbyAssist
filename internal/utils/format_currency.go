package utils

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

func FormatCurrency(input string) string {
	amount := StringToFloat64ElseZero(input)
	withComma := humanize.Commaf(amount)

	return fmt.Sprintf("₱ %s", withComma)
}

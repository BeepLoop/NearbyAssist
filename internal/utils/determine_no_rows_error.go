package utils

import "strings"

func DetermineNoRowsError(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}

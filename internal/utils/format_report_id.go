package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func FormatReportID(id int) string {
	prefix := "REP-"
	return fmt.Sprintf("%s%05d", prefix, id)
}

func ParseReportID(formatted string) (int, error) {
	prefix := "REP-"

	if !strings.HasPrefix(formatted, prefix) {
		return 0, errors.New("invalid report id format")
	}

	numeric := strings.TrimPrefix(formatted, prefix)
	id, err := strconv.Atoi(numeric)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric part in report ID: %s", formatted)
	}

	return id, nil
}

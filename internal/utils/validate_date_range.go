package utils

import (
	"errors"
	"time"
)

const (
    DATE_PARSE_ERROR = "Parsing date error"
)

func ValidateDateRange(start, end string) error {
	now := time.Now().UTC()
	formatDate := "2006-01-02"

	startDate, err := time.Parse(formatDate, start)
	if err != nil {
		return errors.New(DATE_PARSE_ERROR)
	}

	endDate, err := time.Parse(formatDate, end)
	if err != nil {
		return errors.New(DATE_PARSE_ERROR)
	}

	// Validate that the start date is not before the current date
	// Date should not be today or before today
	if startDate.Before(now) || endDate.Before(now) || endDate.Before(startDate) {
		return errors.New("Invalid date")
	}

	return nil
}

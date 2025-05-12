package utils

import (
	"errors"
	"time"
)

const (
	DATE_PARSE_ERR   = "Error parsing date"
	INVALID_DATE_ERR = "Invalid date range"
)

func ValidateDateRange(start, end time.Time) error {
	if start.After(end) || end.Before(start) {
		return errors.New(INVALID_DATE_ERR)
	}

	return nil
}

func ValidateDate(date string) error {
	now := time.Now().UTC()

	input, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		return err
	}

	if input.Before(now) {
		return errors.New(INVALID_DATE_ERR)
	}

	return nil
}

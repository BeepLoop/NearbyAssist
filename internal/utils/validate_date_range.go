package utils

import (
	"errors"
	"fmt"
	"time"
)

const (
	DATE_PARSE_ERR   = "Error parsing date"
	INVALID_DATE_ERR = "Invalid date"
)

func ValidateDateRange(start, end string) error {
	fmt.Println("start date: ", start)
	fmt.Println("end date: ", end)

	now := time.Now().UTC()
	formatDate := time.RFC1123Z

	startDate, err := time.Parse(formatDate, start)
	if err != nil {
		return errors.New(DATE_PARSE_ERR)
	}

	endDate, err := time.Parse(formatDate, end)
	if err != nil {
		return errors.New(DATE_PARSE_ERR)
	}

	// Validate that the start date is not before the current date
	// Date should not be today or before today
	if startDate.Before(now) || endDate.Before(now) || endDate.Before(startDate) {
		return errors.New(INVALID_DATE_ERR)
	}

	return nil
}

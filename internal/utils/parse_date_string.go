package utils

import (
	"fmt"
	"time"
)

func ParseDateString(date string) (time.Time, error) {
	layout := "2006-01-02"

	parsed, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println("error parsing date: ", err.Error())
		return time.Time{}, err
	}

	return parsed, nil
}

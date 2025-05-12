package utils

import (
	"fmt"
	"time"
)

func FormatDate(date string) string {
	if date == "" {
		return ""
	}

	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)
	if err != nil {
		return date
	}

	return t.Format(time.RFC1123)
}

func FormatDateTime(date time.Time) string {
	layout := "2006-01-02 15:04:05"
	return date.Format(layout)
}

func FormatDMY(date string) string {
	if date == "" {
		return ""
	}

	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)
	if err != nil {
		return date
	}

	return t.Format("02 Jan 2006")
}

func DateMonth(date string) string {
	if date == "" {
		return ""
	}

	t, err := time.Parse(time.RFC1123, date)
	if err != nil {
		return date
	}

	return fmt.Sprintf("%s %d", t.Month(), t.Day())
}

func CurrentMonthYear() string {
	format := "January 2006"

	location, _ := time.LoadLocation("Asia/Manila")
	now := time.Now().In(location)

	return now.Format(format)
}

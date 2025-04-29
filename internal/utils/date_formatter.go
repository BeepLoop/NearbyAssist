package utils

import "time"

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

func DateMonth(date string) string {
	if date == "" {
		return ""
	}

	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)
	if err != nil {
		return date
	}

	return t.Format("January 2")
}

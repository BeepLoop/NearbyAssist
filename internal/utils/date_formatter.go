package utils

import "time"

func FormatDate(date string) string {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)
	if err != nil {
		return date
	}

	return t.Format(time.RFC1123)
}

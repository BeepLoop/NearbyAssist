package utils

import "time"

func StringToDateTime(timestamp string) (time.Time, error) {
	layout := "2006-01-02T15:04:05Z"

	return time.Parse(layout, timestamp)
}

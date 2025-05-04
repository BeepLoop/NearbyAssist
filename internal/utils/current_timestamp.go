package utils

import (
	"fmt"
	"time"
)

func CurrentTimeStamp() string {
	format := "2006-01-02T15:04:05Z"

	location, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		fmt.Println("Error timezone: ", err.Error())
		return time.Now().Format(format)
	}

	return time.Now().In(location).Format(format)
}

func CurrentTimeStampNonUTC() string {
	format := "2006-01-02 15:04:05Z"

	location, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		fmt.Println("Error timezone: ", err.Error())
		return time.Now().Format(format)
	}

	return time.Now().In(location).Format(format)
}

package utils

import (
	"log"
	"regexp"
	"strings"
)

// example point: POINT(7.255 125.891)
func LatlongExtractor(point string) []string {
	regex, err := regexp.Compile(`\((.*?)\)`)
	if err != nil {
		log.Fatal("err initializing regex: ", err)
	}

	results := regex.FindStringSubmatch(point)
	if results == nil {
		return nil
	}

	points := strings.Split(results[1], " ")

	return points
}

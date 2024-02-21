package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

// example point: POINT(7.255 125.891)
func LatlongExtractor(point string) (float64, float64) {
	regex, err := regexp.Compile(`\((.*?)\)`)
	if err != nil {
		log.Fatal("err initializing regex: ", err)
	}

	results := regex.FindStringSubmatch(point)
	if results == nil {
		return 0, 0
	}

	points := strings.Split(results[1], " ")

	lat, err := strconv.ParseFloat(points[0], 64)
	if err != nil {
		return 0, 0
	}

	lon, err := strconv.ParseFloat(points[1], 64)
	if err != nil {
		return 0, 0
	}

	return lat, lon
}

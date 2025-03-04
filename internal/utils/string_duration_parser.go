package utils

import (
	"errors"
	"strconv"
	"time"
	"unicode"
)

func ParseStringDuration(duration string) (time.Duration, error) {
	var amount time.Duration
	suffix := ""

	for i, c := range duration {
		if unicode.IsLetter(c) {
			suffix = string(c)
			d := duration[0:i]

			v, err := strconv.Atoi(d)
			if err != nil {
				return time.Second, err
			}

			amount = time.Duration(v)

			break
		}
	}

	switch suffix {
	case "s":
		return time.Second * amount, nil
	case "m":
		return time.Minute * amount, nil
	case "h":
		return time.Hour * amount, nil
	case "d":
		return (time.Hour * 24) * amount, nil
	default:
		return time.Second, errors.New("unknown duration")
	}
}

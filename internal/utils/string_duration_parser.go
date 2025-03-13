package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func FormatDurationToString(d time.Duration) string {
	d = d.Abs()

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := []string{}

	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}

	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}

	if minutes > 0 && len(parts) < 2 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	if seconds > 0 && len(parts) < 2 {
		if seconds == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}

	if len(parts) == 0 {
		return "0 seconds"
	}

	return strings.Join(parts, ", ")
}

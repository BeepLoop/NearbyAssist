package utils

import "strings"

func ParseQuery(query string) map[string]string {
	queries := strings.Split(query, "&")

	params := make(map[string]string)
	for _, q := range queries {
		param := strings.Split(q, "=")
		if len(param) == 2 {
			params[param[0]] = param[1]
		}
	}

	return params
}

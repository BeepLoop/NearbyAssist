package utils

func RemoveStringFromSlice(slice []string, target string) []string {
	var result []string

	for _, item := range slice {
		if item != target {
			result = append(result, item)
		}
	}

	return result
}

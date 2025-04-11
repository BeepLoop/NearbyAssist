package utils

func StringOrNil(input string) *string {
	if input == "" {
		return nil
	}

	return &input
}

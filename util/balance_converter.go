package util

import "strconv"

func ConvertToInt(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic("failed to convert string to int: " + err.Error())
	}
	return intValue, nil
}

package util

import "strconv"

func ConvertToInt(value string) (int, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return int(floatValue * 100), nil
}

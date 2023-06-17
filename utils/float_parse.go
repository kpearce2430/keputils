package utils

import (
	"strconv"
	"strings"
)

// FloatParse common float parser
func FloatParse(inputString string) (float64, error) {

	t := strings.Replace(inputString, ",", "", -1)
	t = strings.Replace(t, "$", "", -1)
	t = strings.Replace(t, "#", "", -1)
	t = strings.Replace(t, "%", "", -1)

	switch t {
	case "N/A":
		return 0.00, nil

	case "":
		return 0.00, nil

	case "Add":
		return 0.00, nil

	default:
		value, err := strconv.ParseFloat(t, 64)
		if err != nil {
			value = 0.00
		}

		return value, err
	}
}

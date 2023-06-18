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

// AsciiString returns a string of only ascii characters.
func AsciiString(str string) string {

	byteString := []byte(str)
	newByte := []byte("")
	for i := 0; i < len(byteString); i++ {
		if byteString[i] >= 32 && byteString[i] <= 127 {
			newByte = append(newByte, byteString[i])
		}
	}

	return string(newByte)
}

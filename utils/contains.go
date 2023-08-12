package utils

import "github.com/sirupsen/logrus"

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			logrus.Debug("contains ", str, " in ", s)
			return true
		}
	}
	return false
}

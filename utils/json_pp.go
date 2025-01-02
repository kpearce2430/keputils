package utils

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func PrettyPrintJson(buf []byte) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, buf, "", "\t"); err != nil {
		logrus.Error("Error:", err.Error())
		return "", err
	}
	return string(prettyJSON.Bytes()), nil
}

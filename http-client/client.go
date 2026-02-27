package http_client

import (
	"net/http"
	"time"
)

//func GetDefaultClient(timeout int64, devMode bool) *req.Client {
//	client := req.C().SetTimeout(time.Duration(timeout) * time.Second)
//	if devMode {
//		client.DevMode()
//	}
//	return client
//}

func GetDefaultClient(timeout int64) *http.Client {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return client
}

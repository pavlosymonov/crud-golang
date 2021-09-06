package http

import (
	"encoding/json"
	"log"
)

var (
	StatusSuccess = "SUCCESS"
	StatusError   = "ERROR"
	Success       = `{"status": "` + StatusSuccess + `"}`
)

type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewResponse(status string, message string) string {
	data, err := json.Marshal(Response{
		Status:  status,
		Message: message,
	})
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(data)
}

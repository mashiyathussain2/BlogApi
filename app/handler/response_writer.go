package handler

import (
	"encoding/json"
	"net/http"

	"blog/app/model"
)

// ResponseWriter will write result in http.ResponseWriter
func ResponseWriter(res http.ResponseWriter, statusCode int, message string, data interface{}) error {
	res.WriteHeader(statusCode)
	httpResponse := model.NewResponse(statusCode, message, data)
	err := json.NewEncoder(res).Encode(httpResponse)
	return err
}

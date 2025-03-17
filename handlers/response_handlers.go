package handlers

import (
	"encoding/json"
	"iv_project/dto"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: statusCode, Message: message, Data: data})
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResult{Code: statusCode, Message: message})
}

package handlers

import (
	"encoding/json"
	"iv_project/dto"
	"net/http"
)

// successResponse sends a standardized success response with a status code, message, and data
func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: statusCode, Message: message, Data: data})
}

// errorResponse sends a standardized error response with a status code and message
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResult{Code: statusCode, Message: message})
}

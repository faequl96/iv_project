package dto

// SuccessResult is used for successful API responses
type SuccessResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResult is used for error API responses
type ErrorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

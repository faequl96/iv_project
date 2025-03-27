package dto

type SuccessResult struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

type ErrorResult struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

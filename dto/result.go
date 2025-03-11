package dto

type SuccessResult struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

type ErrorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

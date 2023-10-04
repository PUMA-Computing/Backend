package utils

type ErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

package utils

type ErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type UnauthorizedError struct {
	Message string `json:"message"`
}

func (u UnauthorizedError) Error() string {
	return u.Message
}

package utils

type CustomError struct {
	ErrorResponse ErrorResponse
}

type ConflictError struct {
	Message string
}

func (ce *ConflictError) Error() string {
	return ce.Message
}

func (ce *CustomError) Error() string {
	return ce.ErrorResponse.Errors[0].Message
}

package utils

type CustomError struct {
	ErrorResponse ErrorResponse
}

func (ce *CustomError) Error() string {
	return ce.ErrorResponse.Errors[0].Message
}

type ConflictError struct {
	Message string
}

func (ce *ConflictError) Error() string {
	return ce.Message
}

type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

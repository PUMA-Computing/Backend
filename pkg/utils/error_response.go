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

type BadRequestError struct {
	Message string `json:"message"`
}

type MaxRegistrationReachedError struct {
	EventID int `json:"event_id"`
}

type AlreadyRegisteredError struct {
	EventID int `json:"event_id"`
}

func (m MaxRegistrationReachedError) Error() string {
	return "Maximum registration limit reached for event with ID: " + string(rune(m.EventID))
}

func (u UnauthorizedError) Error() string {
	return u.Message
}

func (b BadRequestError) Error() string {
	return b.Message
}

func (a AlreadyRegisteredError) Error() string {
	return "User is already registered for event with ID: " + string(rune(a.EventID))
}

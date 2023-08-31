package validation

import (
	"Backend/internal/app/domain/event"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateEvent(event *event.Event) error {
	return validate.Struct(event)
}

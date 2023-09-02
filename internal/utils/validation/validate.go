package validation

import (
	"Backend/internal/app/domain/event"
	"Backend/internal/app/domain/news"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateEvent(event *event.Events) error {
	return validate.Struct(event)
}

func ValidateNews(news *news.News) error {
	return validate.Struct(news)
}

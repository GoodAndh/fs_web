package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotFound error = errors.New("requested item not found")
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		Validator *validator.Validate
	}

	GlobalResponseError struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

)

func (x *XValidator) Validate(stx interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}
	errs := x.Validator.Struct(stx)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			if str, ok := err.Value().(string); ok {
				if str == "" {
					elem.Value = "empty value"
				}
			} else if !ok {
				elem.Value = err.Value()
			}
			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Error = true
			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}

func WriteJson(c *fiber.Ctx, status int, message string, v interface{}) error {
	return c.Status(status).JSON(&GlobalResponseError{
		Status:  status,
		Message: message,
		Data:    v,
	})
}

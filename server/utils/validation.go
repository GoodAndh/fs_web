package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotFound error = errors.New("requested item not found")
)

type (
	ErrorResponseValidate struct {
		Error       bool
		FailedField string
		Tag         string
		Message     string
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

func (x *XValidator) validate(stx interface{}) []ErrorResponseValidate {
	validationErrors := []ErrorResponseValidate{}
	errs := x.Validator.Struct(stx)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponseValidate
			switch err.Tag() {
			case "required":
				elem.Message = fmt.Sprintf("form %v wajib di isi", err.Field())
			case "email":
				elem.Message = fmt.Sprintf("form %v tidak sesuai format email", err.Field())
			case "eqfield":
				elem.Message = fmt.Sprintf("form %v harus sama dengan form %v", err.Field(), err.Param())
			case "oneof":
				elem.Message = fmt.Sprintf("only accepted one of = %v", err.Param())
			case "min":
				elem.Message = fmt.Sprintf("minimal di isi dengan %v karakter", err.Param())
			}
			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Error = true
			elem.Value = err.Value()
			validationErrors = append(validationErrors, elem)

		}
	}
	return validationErrors
}

func (x *XValidator) Validate(stx interface{}) map[string]any {

	errMsg := make(map[string]any, 0)
	if errs := x.validate(stx); len(errs) > 0 && errs[0].Error {
		for _, err := range errs {
			errMsg[err.FailedField] = err.Message
		}
	}
	return errMsg
}

func WriteJson(c *fiber.Ctx, status int, message string, v interface{}) error {
	return c.Status(status).JSON(&GlobalResponseError{
		Status:  status,
		Message: message,
		Data:    v,
	})
}

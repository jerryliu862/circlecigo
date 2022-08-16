package util

import (
	"github.com/go-playground/validator"
)

func ValidateData(data interface{}) error {
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		validationErrors, _ := err.(validator.ValidationErrors)
		return validationErrors
	}
	return nil
}

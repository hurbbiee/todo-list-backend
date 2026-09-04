package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	registerCustomValidations()
}

func registerCustomValidations() {
	validate.RegisterValidation("date_ymd", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		_, err := time.Parse("2006-01-02", dateStr)
		return err == nil
	})
}

func Validate(i interface{}) error {
	return validate.Struct(i)
}

func ValidateStruct(i interface{}) map[string]string {
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = e.Tag()
	}

	return errors
}

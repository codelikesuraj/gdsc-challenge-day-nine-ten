package models

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErr map[string]string

func GetValidationErrs(ve validator.ValidationErrors) []ValidationErr {
	var errs []ValidationErr

	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		tag := fe.Tag()
		param := fe.Param()

		switch tag {
		case "required":
			errs = append(errs, ValidationErr{field: field + " is required"})
		case "min":
			errs = append(errs, ValidationErr{field: field + " should be at least " + param + " characters "})
		case "max":
			errs = append(errs, ValidationErr{field: field + " should be at least " + param + " characters "})
		default:
			errs = append(errs, ValidationErr{"unknown": "unknown error"})
		}
	}
	return errs
}

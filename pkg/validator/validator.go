package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateBody(input interface{}) map[string]string
}

type validatorImpl struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return &validatorImpl{
		validate: validator.New(),
	}
}

func (v *validatorImpl) ValidateBody(input interface{}) map[string]string {
	if err := v.validate.Struct(input); err != nil {
		res := map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			tag := err.Tag()
			switch tag {
			case "required":
				res[field] = fmt.Sprintf("%s is required", field)
			case "email":
				res[field] = fmt.Sprintf("%s must be a valid email", field)
			case "min":
				res[field] = fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
			default:
				res[field] = fmt.Sprintf("%s is not valid", field)
			}
		}

		return res
	}
	return nil
}

package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type validationError struct {
	validator.FieldError
}

type ValidationError []string

type ValidationFieldLevel validator.FieldLevel

type ValidationFieldFunc func(ValidationFieldLevel) bool

type ValidationField struct {
	Tag string
	Fn  ValidationFieldFunc
}

var validate *validator.Validate

var registeredVfs []ValidationField = []ValidationField{}

func Instantiate() error {
	validate = validator.New()
	for _, vf := range registeredVfs {
		if err := validate.RegisterValidation(vf.Tag, vf.Fn.toValidatorFunc); err != nil {
			return err
		}
	}
	return nil
}

func ValidateStruct(i interface{}) ValidationError {
	errs := validate.Struct(i)
	if errs == nil {
		return nil
	}

	var vErrs []validator.FieldError

	for _, err := range errs.(validator.ValidationErrors) {
		ve := validationError{err}
		vErrs = append(vErrs, ve)
	}

	return messages(vErrs)
}

func messageF(v validator.FieldError) string {
	return fmt.Sprintf(
		"Field '%s', Cause: '%s'",
		v.Field(),
		v.Tag(),
	)
}

func messages(v []validator.FieldError) ValidationError {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, messageF(err))
	}
	return errs
}

func RegisterField(vf ValidationField) {
	registeredVfs = append(registeredVfs, vf)
}

func (vfn ValidationFieldFunc) toValidatorFunc(fl validator.FieldLevel) bool {
	return vfn(fl)
}

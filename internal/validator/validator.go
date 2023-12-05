package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	validator.FieldError
}

type ValidationErrors []ValidationError

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

func (vfn ValidationFieldFunc) toValidatorFunc(fl validator.FieldLevel) bool {
	return vfn(fl)
}

func (v ValidationError) String() string {
	return fmt.Sprintf(
		"Field '%s', Cause: '%s'",
		// v.StructField(),
		v.Field(),
		v.Tag(),
	)
}

func (v ValidationErrors) String() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.String())
	}
	return errs
}

func RegisterField(vf ValidationField) {
	registeredVfs = append(registeredVfs, vf)
}

func Validate(i interface{}) ValidationErrors {
	errs := validate.Struct(i)
	if errs == nil {
		return nil
	}

	var vErrs []ValidationError

	for _, err := range errs.(validator.ValidationErrors) {
		ve := ValidationError{err}
		vErrs = append(vErrs, ve)
	}

	return vErrs
}

// func ValidateMap(data map[string]interface{}, rules ValidationRules) ValidationErrors {
// 	logger.Info(rules)
// 	logger.Info(data)
// 	errs := ValidationErrors{}
// 	for key, val := range data {
// 		if rv, ok := rules[key]; ok {
// 			//key can be string or map shoulda check
// 			//since this approach doesnt show any error messages it looks very silly
// 			//decided to not use it at all
// 			logger.Debug("key: ", key, " val: ", val, " rule: ", rv)
// 			if e := validate.Var(val, rv).(validator.ValidationErrors); e != nil {
// 				for _, v := range e {
// 					errs = append(errs, ValidationError{v})
// 				}
// 			}
// 		}
// 	}
// 	return errs
// }

package validator

import (
	"reflect"
	"strings"
	"sync"

	lib "github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(i any) *Error
}

type validator struct {
	validate *lib.Validate
}

var (
	si              *validator // since validator caches structs
	instanciateOnce sync.Once  // gotta instanicate once
)

func Instance() Validator {
	instanciateOnce.Do(func() {
		va := lib.New(lib.WithRequiredStructEnabled())

		va.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			// skip if tag key says it should be ignored
			if name == "-" {
				return ""
			}
			return name
		})

		i := &validator{
			validate: va,
		}
		si = i
	})
	return si
}

// Validates a struct
func (v *validator) Validate(i any) *Error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	errs := []FieldError{}
	for _, e := range err.(lib.ValidationErrors) {
		errs = append(errs, FieldError{
			Field: e.Field(),

			Value: e.Value(),
			Cause: e.Tag(),
		})
	}

	return &Error{
		Fields: errs,
	}

}

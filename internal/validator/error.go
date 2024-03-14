package validator

type Error struct {
	Fields []FieldError `json:"errors"`
}

type FieldError struct {
	Field string `json:"f"`
	Cause string `json:"c"`
	Value any    `json:"v"`
}

func (fe *FieldError) ToMap() map[string]interface{} {
	return map[string]interface{}{
		fe.Field: map[string]interface{}{
			"c:": fe.Cause,
			"v":  fe.Value,
		},
	}
}

func (e *Error) ToMap() []map[string]interface{} {
	errs := []map[string]interface{}{}
	for _, fe := range e.Fields {
		errs = append(errs, fe.ToMap())
	}
	return errs
}

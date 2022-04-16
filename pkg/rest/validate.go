package rest

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Validation struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type CanValidate interface {
	Validate(v Validator)
}

type ValidateContext map[string]any

type Validator struct {
	Path        []string
	Validations *[]Validation
	Context     ValidateContext
}

var _ ReflectConsumer = Validator{}

func NewValidator() Validator {
	return Validator{
		Validations: &[]Validation{},
		Context:     ValidateContext{},
	}
}

func (v Validator) Field(field string) Validator {
	return Validator{
		Path:        append(v.Path, field),
		Validations: v.Validations,
		Context:     v.Context,
	}
}

func (v *Validator) Add(format string, args ...any) {
	*v.Validations = append(*v.Validations, Validation{
		Path:    strings.Join(v.Path, "."),
		Message: fmt.Sprintf(format, args...),
	})
}

func (v *Validator) AddField(field string, format string, args ...any) {
	*v.Validations = append(*v.Validations, Validation{
		Path:    strings.Join(append(v.Path, field), "."),
		Message: fmt.Sprintf(format, args...),
	})
}

func (v Validator) Validate(value any) {
	ValidateReflector.Consume(value, v)
}

func (v Validator) IsValid() bool {
	return len(*v.Validations) == 0
}

func (v Validator) Consume(rv reflect.Value) {
	if can, ok := rv.Interface().(CanValidate); ok {
		can.Validate(v)
	}
}

func (v Validator) ForIndex(index int) ReflectConsumer {
	return v.Field(strconv.Itoa(index))
}

func (v Validator) ForField(field reflect.StructField) ReflectConsumer {
	fieldName := field.Name
	if jsonName, ok := field.Tag.Lookup("json"); ok {
		fieldName = jsonName
	}
	if validate, ok := field.Tag.Lookup("validate"); ok {
		if strings.Index(validate, ",skip") != -1 {
			return nil
		}
		if strings.Index(validate, "-") != -1 {
			return v
		}
		split := strings.Split(validate, ",")
		if len(split[0]) > 0 {
			fieldName = split[0]
		}
	}

	return v.Field(fieldName)
}

func (v Validator) ForKey(key string) ReflectConsumer {
	return v.Field(key)
}

var (
	CanValidateType   = reflect.TypeOf((*CanValidate)(nil)).Elem()
	ValidateReflector = NewReflector(func(t reflect.Type) bool {
		return t.Implements(CanValidateType)
	})
)

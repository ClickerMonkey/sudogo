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

func NewValidator() Validator {
	return Validator{
		Validations: &[]Validation{},
		Context:     ValidateContext{},
	}
}

func (v Validator) Item(index int) Validator {
	return v.Field(strconv.Itoa(index))
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
	reflectValue := reflect.ValueOf(value)
	validateFuncs := GetValidators(reflectValue.Type())
	doFuncsValidate(validateFuncs, reflectValue, v)
}

func (v Validator) IsValid() bool {
	return len(*v.Validations) == 0
}

type ValidateFunc func(value reflect.Value, validator Validator)

var (
	ValidatorMap    = map[reflect.Type][]ValidateFunc{}
	CanValidateType = reflect.TypeOf((*CanValidate)(nil)).Elem()
)

func doCanValidate(value reflect.Value, validator Validator) {
	if can, ok := value.Interface().(CanValidate); ok {
		can.Validate(validator)
	}
}

func doFuncsValidate(funcs []ValidateFunc, value reflect.Value, validator Validator) {
	for _, f := range funcs {
		f(value, validator)
	}
}

func doPointerValidate(funcs []ValidateFunc) ValidateFunc {
	if len(funcs) == 0 {
		return nil
	}

	return func(value reflect.Value, validator Validator) {
		if !value.IsNil() {
			doFuncsValidate(funcs, value.Elem(), validator)
		}
	}
}

func doInterfaceValidate(funcs []ValidateFunc) ValidateFunc {
	if len(funcs) == 0 {
		return nil
	}

	return func(value reflect.Value, validator Validator) {
		doFuncsValidate(funcs, value.Elem(), validator)
	}
}

func doFieldValidate(parent reflect.Type, fieldIndex int) ValidateFunc {
	field := parent.Field(fieldIndex)
	funcs := GetValidators(field.Type)
	fieldName := field.Tag.Get("json")
	if fieldName == "" {
		fieldName = field.Name
	}

	if len(funcs) == 0 {
		return nil
	}

	return func(value reflect.Value, validator Validator) {
		fieldValue := value.Field(fieldIndex)
		if field.Anonymous {
			doFuncsValidate(funcs, fieldValue, validator)
		} else {
			doFuncsValidate(funcs, fieldValue, validator.Field(fieldName))
		}
	}
}

func doSliceValidate(funcs []ValidateFunc) ValidateFunc {
	if len(funcs) == 0 {
		return nil
	}

	return func(value reflect.Value, validator Validator) {
		for i := 0; i < value.Len(); i++ {
			doFuncsValidate(funcs, value.Index(i), validator.Item(i))
		}
	}
}

func doMapValidate(funcs []ValidateFunc) ValidateFunc {
	if len(funcs) == 0 {
		return nil
	}

	return func(value reflect.Value, validator Validator) {
		iter := value.MapRange()
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()

			doFuncsValidate(funcs, value, validator.Field(key.String()))
		}
	}
}

func AddValidator(typ reflect.Type, f ValidateFunc) {
	if f != nil {
		ValidatorMap[typ] = append(ValidatorMap[typ], f)
	}
}

func GetValidators(typ reflect.Type) []ValidateFunc {
	if funcs, exists := ValidatorMap[typ]; exists {
		return funcs
	}

	ValidatorMap[typ] = []ValidateFunc{}

	if typ.Implements(CanValidateType) {
		AddValidator(typ, doCanValidate)
	}

	switch typ.Kind() {
	case reflect.Pointer:
		AddValidator(typ, doPointerValidate(GetValidators(typ.Elem())))
	case reflect.Interface:
		AddValidator(typ, doInterfaceValidate(GetValidators(typ.Elem())))
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			AddValidator(typ, doFieldValidate(typ, i))
		}
	case reflect.Array, reflect.Slice:
		AddValidator(typ, doSliceValidate(GetValidators(typ.Elem())))
	case reflect.Map:
		AddValidator(typ, doMapValidate(GetValidators(typ.Elem())))
	}

	return ValidatorMap[typ]
}

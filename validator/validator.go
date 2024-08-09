package validator

import (
	"fmt"
)

type Validator interface {
	Uint(n uint, field string) error
	Uint8(n uint8, field string) error
	String(s string, field string) error
	Nil(value interface{}, field string) error
}

type validator struct {
}

func New() Validator {
	return &validator{}
}

func (validator) Nil(value interface{}, field string) error {
	if value == nil {
		return fmt.Errorf("campo nulo '%s' é requerido", field)
	}
	return nil
}

func (v *validator) String(param string, field string) error {
	if len(param) == 0 {
		return v.err("string", field)
	}

	return nil
}

func (v *validator) Uint(ID uint, field string) error {
	if ID == 0 {
		return v.err("uint", field)
	}

	return nil
}

func (v *validator) Uint8(ID uint8, field string) error {
	if ID == 0 {
		return v.err("uint8", field)
	}

	return nil
}

func (validator) err(_type, field string) error {
	return fmt.Errorf("campo '%s' do tipo '%s' é requerido", field, _type)
}

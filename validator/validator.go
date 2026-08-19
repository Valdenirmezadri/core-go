package validator

import (
	"fmt"
	"net/url"
	"time"
)

type Validator interface {
	Uint(n uint, field string) error
	Uint8(n uint8, field string) error
	String(s string, field string) error
	Bytes(data []byte, field string) error
	Nil(value any, field string) error
	Time(value time.Time, field string) error
	TimeBetween(start, end time.Time) error
	URL(url, field string) error
}

type validator struct {
}

func New() Validator {
	return &validator{}
}

func (validator) Nil(value any, field string) error {
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

func (v *validator) Bytes(data []byte, field string) error {
	if len(data) == 0 {
		return v.err("[]byte", field)
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

func (v *validator) URL(urlStr, field string) error {
	if err := v.String(urlStr, field); err != nil {
		return err
	}

	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return fmt.Errorf("%s: %w", v.err("string url é inválido,", field), err)
	}

	return nil
}

func (validator) err(_type, field string) error {
	return fmt.Errorf("campo '%s' do tipo '%s' é requerido", field, _type)
}

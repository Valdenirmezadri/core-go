package validator

import (
	"embed"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"time"

	"github.com/Valdenirmezadri/core-go/context/corectx"
	"github.com/Valdenirmezadri/core-go/i18n"
)

//go:embed locales/*.toml
var localeFS embed.FS

type Validator interface {
	Uint(ctx corectx.Context, n uint, field string) error
	Uint8(ctx corectx.Context, n uint8, field string) error
	String(ctx corectx.Context, s string, field string) error
	Bytes(ctx corectx.Context, data []byte, field string) error
	Nil(ctx corectx.Context, value any, field string) error
	Slice(ctx corectx.Context, length int, field string) error
	Time(ctx corectx.Context, value time.Time, field string) error
	TimeBetween(ctx corectx.Context, start, end time.Time) error
	URL(ctx corectx.Context, url, field string) error
	Email(ctx corectx.Context, email, field string) error
	Password(ctx corectx.Context, password, passwordMatch string) error
}


type validator struct {
	t i18n.Switch
}

func New(logger func() i18n.Logger) (Validator, error) {
	t, err := i18n.New(localeFS, logger)
	if err != nil {
		return nil, err
	}

	return &validator{t: t}, nil
}

func (v *validator) Slice(ctx corectx.Context, length int, field string) error {
	if length == 0 {
		return v.err(ctx, "slice", field)
	}
	return nil
}

func (v *validator) Nil(ctx corectx.Context, value any, field string) error {
	if value == nil {
		return errors.New(v.t.S(ctx, "NilFieldRequired", map[string]string{"Field": field}))
	}
	return nil
}

func (v *validator) String(ctx corectx.Context, param string, field string) error {
	if len(param) == 0 {
		return v.err(ctx, "string", field)
	}
	return nil
}

func (v *validator) Bytes(ctx corectx.Context, data []byte, field string) error {
	if len(data) == 0 {
		return v.err(ctx, "[]byte", field)
	}
	return nil
}

func (v *validator) Uint(ctx corectx.Context, ID uint, field string) error {
	if ID == 0 {
		return v.err(ctx, "uint", field)
	}
	return nil
}

func (v *validator) Uint8(ctx corectx.Context, ID uint8, field string) error {
	if ID == 0 {
		return v.err(ctx, "uint8", field)
	}
	return nil
}

func (v *validator) URL(ctx corectx.Context, urlStr, field string) error {
	if err := v.String(ctx, urlStr, field); err != nil {
		return err
	}

	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return fmt.Errorf("%s: %w",
			v.t.S(ctx, "InvalidURL", map[string]string{"Field": field}), err)
	}

	return nil
}

func (v *validator) Email(ctx corectx.Context, email, field string) error {
	if err := v.String(ctx, email, field); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("%s: %w",
			v.t.S(ctx, "InvalidEmail", map[string]string{"Field": field}), err)
	}

	return nil
}

func (v *validator) Password(ctx corectx.Context, password, passwordMatch string) error {
	if password == "" {
		return errors.New(v.t.S(ctx, "PasswordEmpty", nil))
	}

	if password != passwordMatch {
		return errors.New(v.t.S(ctx, "PasswordMismatch", nil))
	}

	return nil
}

func (v *validator) err(ctx corectx.Context, _type, field string) error {
	return errors.New(v.t.S(ctx, "FieldTypeRequired", map[string]string{
		"Field": field,
		"Type":  _type,
	}))
}

package appctx

import (
	helperecho "github.com/Valdenirmezadri/core-go/helpers/http/echo"
	"github.com/Valdenirmezadri/core-go/hterr"
	secure "github.com/Valdenirmezadri/core-go/jwt"
	"github.com/Valdenirmezadri/core-go/text"
	"github.com/Valdenirmezadri/core-go/validator"
)

type Tools interface {
	JWT() secure.SecureHandler
	Err() hterr.Err
	Validator() validator.Validator
	Text() text.Helper
	HTTP() helperecho.Helper
}

type tools struct {
	jwt       secure.SecureHandler
	err       hterr.Err
	validator validator.Validator
	text      text.Helper
	http      helperecho.Helper
}

func NewTools(jwt secure.SecureHandler, err hterr.Err, v validator.Validator, t text.Helper, h helperecho.Helper) Tools {
	return &tools{
		jwt:       jwt,
		err:       err,
		validator: v,
		text:      t,
		http:      h,
	}
}

func (t *tools) JWT() secure.SecureHandler      { return t.jwt }
func (t *tools) Err() hterr.Err                 { return t.err }
func (t *tools) Validator() validator.Validator { return t.validator }
func (t *tools) Text() text.Helper              { return t.text }
func (t *tools) HTTP() helperecho.Helper        { return t.http }

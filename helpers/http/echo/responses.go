package helperecho

import (
	"embed"
	"errors"
	"net/http"
	"time"

	"github.com/Valdenirmezadri/core-go/context/corectx"
	"github.com/Valdenirmezadri/core-go/i18n"
	"github.com/Valdenirmezadri/core-go/slices"
	"github.com/labstack/echo/v4"
)

//go:embed locales/*.toml
var localeFS embed.FS

type Helper interface {
	Bind(c echo.Context, i interface{}) error

	RequiredPropToString(ctx corectx.Context, c echo.Context, propName string) (string, error)
	PropToString(ctx corectx.Context, c echo.Context, propName string, required bool) (string, error)
	RequiredPropToUint(ctx corectx.Context, c echo.Context, propName string) (uint, error)
	PropToUint(ctx corectx.Context, c echo.Context, propName string, required bool) (uint, error)
	RequiredPropToBool(ctx corectx.Context, c echo.Context, propName string) (bool, error)
	PropToBool(ctx corectx.Context, c echo.Context, propName string, required bool) (bool, error)

	RequiredFormToString(ctx corectx.Context, c echo.Context, param string) (string, error)
	RequiredFormToUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	FormToString(ctx corectx.Context, c echo.Context, param string) (string, error)
	FormToUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	FormBool(ctx corectx.Context, c echo.Context, param string) (bool, error)
	FormToBool(ctx corectx.Context, c echo.Context, param string) (bool, error)

	RequiredParamToString(ctx corectx.Context, c echo.Context, param string) (string, error)
	RequiredParamToUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	ParamToString(ctx corectx.Context, c echo.Context, param string) (string, error)
	ParamToUint(ctx corectx.Context, c echo.Context, param string) (uint, error)

	RequiredQueryToString(ctx corectx.Context, c echo.Context, param string) (string, error)
	RequiredQueryToUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	QueryToString(ctx corectx.Context, c echo.Context, param string) (string, error)
	QueryToUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	QueryToBool(ctx corectx.Context, c echo.Context, param string) (bool, error)

	RequiredQueryToDateTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	QueryToDateTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	RequiredQueryToDate(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	QueryToDate(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	RequiredQueryToTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	QueryToTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)

	Response(c echo.Context, data any) error
	ResponseMessage(c echo.Context, message string) error
	ResponseAction(c echo.Context, message string, data any) error
	ResponseCodeErr(c echo.Context, code int, errs ...error) error
	ResponseErr(c echo.Context, err error) error
	BadRequestErr(c echo.Context, err error) error

	Relay(c echo.Context, code int, b []byte) error
	RelayErr(c echo.Context, code int, b []byte) error
	RelayResponse(c echo.Context, code int, b []byte) error

	// Aliases
	ReqPropStr(ctx corectx.Context, c echo.Context, propName string) (string, error)
	PropStr(ctx corectx.Context, c echo.Context, propName string, required bool) (string, error)
	ReqPropUint(ctx corectx.Context, c echo.Context, propName string) (uint, error)
	PropUint(ctx corectx.Context, c echo.Context, propName string, required bool) (uint, error)
	ReqPropBool(ctx corectx.Context, c echo.Context, propName string) (bool, error)
	PropBool(ctx corectx.Context, c echo.Context, propName string, required bool) (bool, error)

	ReqFormStr(ctx corectx.Context, c echo.Context, param string) (string, error)
	ReqFormUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	FormStr(ctx corectx.Context, c echo.Context, param string) (string, error)
	FormUint(ctx corectx.Context, c echo.Context, param string) (uint, error)

	ReqParamStr(ctx corectx.Context, c echo.Context, param string) (string, error)
	ReqParamUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	ParamStr(ctx corectx.Context, c echo.Context, param string) (string, error)
	ParamUint(ctx corectx.Context, c echo.Context, param string) (uint, error)

	ReqQueryStr(ctx corectx.Context, c echo.Context, param string) (string, error)
	ReqQueryUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	QueryStr(ctx corectx.Context, c echo.Context, param string) (string, error)
	QueryUint(ctx corectx.Context, c echo.Context, param string) (uint, error)
	QueryBool(ctx corectx.Context, c echo.Context, param string) (bool, error)

	ReqQueryDateTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	QueryDateTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	ReqQueryDate(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	QueryDate(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	ReqQueryTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	QueryTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)
	BadReqErr(c echo.Context, err error) error

	WriteSSE(e echo.Context, ID uint, kind string, data any) error
	WriteSSEPong(e echo.Context) error
	Res(c echo.Context, data any) error
	ResMsg(c echo.Context, message string) error
	ResAction(c echo.Context, message string, data any) error
	ResCodeErr(c echo.Context, code int, errs ...error) error
	ResErr(c echo.Context, err error) error
	RelayRes(c echo.Context, code int, b []byte) error
}

// HelperEcho is an alias for backward compatibility
type HelperEcho = Helper

type helper struct {
	t i18n.Switch
}

func New(logger func() i18n.Logger) (Helper, error) {
	t, err := i18n.New(localeFS, logger)
	if err != nil {
		return nil, err
	}

	return &helper{t: t}, nil
}

func (h *helper) ResponseMessage(c echo.Context, message string) error {
	return h.ResponseAction(c, message, nil)
}

func (h *helper) Response(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, Result{
		Data: data,
	})
}

func (h *helper) Relay(c echo.Context, code int, b []byte) error {
	if code != 200 {
		return h.RelayErr(c, code, b)
	}

	return h.RelayResponse(c, code, b)
}

func (h *helper) RelayResponse(c echo.Context, code int, b []byte) error {
	old, err := Result{}.Unmarshall(b)
	if err != nil {
		return err
	}

	return c.JSON(code, Result{
		Message: old.Message,
		Data:    old.Data,
	})
}

func (h *helper) ResponseAction(c echo.Context, message string, data any) error {
	return h.response(c, http.StatusOK, message, data)
}

func (h *helper) response(c echo.Context, code int, message string, data any) error {
	return c.JSON(code, Result{Message: message, Data: data})
}

func (h *helper) BadRequestErr(c echo.Context, err error) error {
	return h.ResponseCodeErr(c, http.StatusBadRequest, err)
}

func (h *helper) ResponseErr(c echo.Context, err error) error {
	return h.ResponseCodeErr(c, http.StatusInternalServerError, err)
}

func (h *helper) RelayErr(c echo.Context, code int, b []byte) error {
	old, err := Err{}.Unmarshall(b)
	if err != nil {
		return err
	}

	errs := slices.Map(old.Errors, func(err string) error { return errors.New(err) })
	return h.ResponseCodeErr(c, code, errs...)
}

func (h *helper) ResponseCodeErr(c echo.Context, code int, errs ...error) error {
	var errsStr []string
	for _, err := range errs {
		errsStr = append(errsStr, err.Error())
	}

	return c.JSON(code, Err{
		Errors: errsStr,
	})
}

// Aliases
func (h *helper) ResMsg(c echo.Context, message string) error { return h.ResponseMessage(c, message) }
func (h *helper) Res(c echo.Context, data any) error          { return h.Response(c, data) }
func (h *helper) RelayRes(c echo.Context, code int, b []byte) error {
	return h.RelayResponse(c, code, b)
}
func (h *helper) ResAction(c echo.Context, message string, data any) error {
	return h.ResponseAction(c, message, data)
}
func (h *helper) BadReqErr(c echo.Context, err error) error { return h.BadRequestErr(c, err) }
func (h *helper) ResErr(c echo.Context, err error) error    { return h.ResponseErr(c, err) }
func (h *helper) ResCodeErr(c echo.Context, code int, errs ...error) error {
	return h.ResponseCodeErr(c, code, errs...)
}

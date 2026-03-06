package helperecho

import (
	"errors"
	"net/http"
	"time"

	"github.com/Valdenirmezadri/core-go/slices"
	"github.com/labstack/echo/v4"
)

type HelperEcho interface {
	Bind(c echo.Context, i interface{}) error

	RequiredPropToString(c echo.Context, propName string) (string, error)
	PropToString(c echo.Context, propName string, required bool) (string, error)
	RequiredPropToUint(c echo.Context, propName string) (uint, error)
	PropToUint(c echo.Context, propName string, required bool) (uint, error)
	RequiredPropToBool(c echo.Context, propName string) (bool, error)
	PropToBool(c echo.Context, propName string, required bool) (bool, error)

	RequiredFormToString(c echo.Context, param string) (string, error)
	RequiredFormToUint(c echo.Context, param string) (uint, error)
	FormToString(c echo.Context, param string) (string, error)
	FormToUint(c echo.Context, param string) (uint, error)

	RequiredParamToString(c echo.Context, param string) (string, error)
	RequiredParamToUint(c echo.Context, param string) (uint, error)
	ParamToString(c echo.Context, param string) (string, error)
	ParamToUint(c echo.Context, param string) (uint, error)

	RequiredQueryToString(c echo.Context, param string) (string, error)
	RequiredQueryToUint(c echo.Context, param string) (uint, error)
	QueryToString(c echo.Context, param string) (string, error)
	QueryToUint(c echo.Context, param string) (uint, error)

	QueryToBool(c echo.Context, param string) (bool, error)

	RequiredQueryToDateTime(c echo.Context, param string) (time.Time, error)
	QueryToDateTime(c echo.Context, param string) (time.Time, error)
	RequiredQueryToDate(c echo.Context, param string) (time.Time, error)
	QueryToDate(c echo.Context, param string) (time.Time, error)
	RequiredQueryToTime(c echo.Context, param string) (time.Time, error)
	QueryToTime(c echo.Context, param string) (time.Time, error)

	Response(c echo.Context, data any) error
	ResponseMessage(c echo.Context, message string) error
	ResponseAction(c echo.Context, message string, data any) error
	ResponseCodeErr(c echo.Context, code int, errs ...error) error
	ResponseErr(c echo.Context, err error) error
	BadRequestErr(c echo.Context, err error) error
}

type helperEcho struct {
}

func New() HelperEcho {
	return &helperEcho{}
}

func (h helperEcho) ResponseMessage(c echo.Context, message string) error {
	return h.ResponseAction(c, message, nil)
}

func (helperEcho) Response(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, Result{
		Data: data,
	})
}

func (r *helperEcho) RelayResponse(c echo.Context, b []byte) error {
	old, err := Result{}.Unmarshall(b)
	if err != nil {
		return err
	}

	return r.ResponseAction(c, old.Message, old.Data)
}

func (helperEcho) ResponseAction(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, Result{
		Message: message,
		Data:    data,
	})

}

func (r *helperEcho) BadRequestErr(c echo.Context, err error) error {
	return r.ResponseCodeErr(c, http.StatusBadRequest, err)
}

func (r *helperEcho) ResponseErr(c echo.Context, err error) error {
	return r.ResponseCodeErr(c, http.StatusInternalServerError, err)
}

func (r *helperEcho) RelayErr(c echo.Context, code int, b []byte) error {
	old, err := Err{}.Unmarshall(b)
	if err != nil {
		return err
	}

	errs := slices.Map(old.Errors, func(err string) error { return errors.New(err) })
	return r.ResponseCodeErr(c, code, errs...)
}

func (helperEcho) ResponseCodeErr(c echo.Context, code int, errs ...error) error {
	var errsStr []string
	for _, err := range errs {
		errsStr = append(errsStr, err.Error())
	}

	return c.JSON(code, Err{
		Errors: errsStr,
	})
}

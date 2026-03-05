package helperecho

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type HelperEcho interface {
	Bind(c echo.Context, i interface{}) error

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

func (helperEcho) ResponseCodeErr(c echo.Context, code int, errs ...error) error {
	var errsStr []string
	for _, err := range errs {
		errsStr = append(errsStr, err.Error())
	}
	return c.JSON(code, Err{
		Errors: errsStr,
	})
}

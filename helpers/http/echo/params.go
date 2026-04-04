package helperecho

import (
	"errors"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *helper) Bind(c echo.Context, i interface{}) error {
	if err := c.Bind(&i); err != nil {
		return err
	}

	return nil
}

func (h *helper) RequiredFormToString(c echo.Context, param string) (string, error) {
	return h.formToString(c, param, true)
}

func (h *helper) FormToString(c echo.Context, param string) (string, error) {
	return h.formToString(c, param, false)
}

func (h *helper) formToString(c echo.Context, param string, required bool) (string, error) {
	str := c.FormValue(param)
	if str == "" && required {
		return "", h.errT("FormRequired", param)
	}

	return str, nil
}

func (h *helper) RequiredFormToUint(c echo.Context, param string) (uint, error) {
	return h.formToUint(c, param, true)
}

func (h *helper) FormToUint(c echo.Context, param string) (uint, error) {
	return h.formToUint(c, param, false)
}

func (h *helper) formToUint(c echo.Context, param string, required bool) (uint, error) {
	strID := c.FormValue(param)
	if strID == "" {
		if required {
			return 0, h.errT("FormRequired", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helper) ParamToString(c echo.Context, param string) (string, error) {
	return h.paramToString(c, param, false)
}

func (h *helper) RequiredParamToString(c echo.Context, param string) (string, error) {
	return h.paramToString(c, param, true)
}

func (h *helper) paramToString(c echo.Context, param string, required bool) (string, error) {
	str := c.Param(param)
	if str == "" && required {
		return "", h.errT("PathRequired", param)
	}

	return str, nil
}

func (h *helper) RequiredParamToUint(c echo.Context, param string) (uint, error) {
	return h.paramToUint(c, param, true)
}
func (h *helper) ParamToUint(c echo.Context, param string) (uint, error) {
	return h.paramToUint(c, param, false)
}

func (h *helper) paramToUint(c echo.Context, param string, required bool) (uint, error) {
	strID := c.Param(param)
	if strID == "" || strID == "0" {
		if required {
			return 0, h.errT("PathRequired", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helper) RequiredQueryToString(c echo.Context, param string) (string, error) {
	return h.queryToString(c, param, true)
}

func (h *helper) QueryToString(c echo.Context, param string) (string, error) {
	return h.queryToString(c, param, false)
}

func (h *helper) queryToString(c echo.Context, param string, required bool) (string, error) {
	str := strings.TrimSpace(c.QueryParam(param))
	if str == "" && required {
		return "", h.errT("QueryRequired", param)
	}

	return str, nil
}

func (h *helper) RequiredQueryToUint(c echo.Context, param string) (uint, error) {
	return h.queryToUint(c, param, true)
}

func (h *helper) QueryToUint(c echo.Context, param string) (uint, error) {
	return h.queryToUint(c, param, false)
}

func (h *helper) queryToUint(c echo.Context, param string, required bool) (uint, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return 0, h.errT("QueryRequired", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helper) FormBool(c echo.Context, param string) (bool, error) {
	str := strings.ToLower(strings.TrimSpace(c.FormValue(param)))
	if str == "" {
		return false, nil
	}

	switch str {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, h.errT("InvalidBoolForm", param)
	}
}

func (h *helper) QueryToBool(c echo.Context, param string) (bool, error) {
	str := strings.ToLower(strings.TrimSpace(c.QueryParam(param)))
	if str == "" {
		return false, nil
	}

	switch str {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, h.errT("InvalidBoolQuery", param)
	}
}

func (h *helper) convStrToUint(str string, required bool) (uint, error) {
	ID, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		if required {
			return 0, err
		}

		return 0, err
	}

	return uint(ID), nil
}

func (h *helper) errT(id, field string) error {
	return errors.New(h.t.S(nil, id, map[string]string{"Field": field}))
}

// Aliases
func (h *helper) ReqFormStr(c echo.Context, param string) (string, error)  { return h.RequiredFormToString(c, param) }
func (h *helper) FormStr(c echo.Context, param string) (string, error)     { return h.FormToString(c, param) }
func (h *helper) ReqFormUint(c echo.Context, param string) (uint, error)   { return h.RequiredFormToUint(c, param) }
func (h *helper) FormUint(c echo.Context, param string) (uint, error)      { return h.FormToUint(c, param) }
func (h *helper) ReqParamStr(c echo.Context, param string) (string, error) { return h.RequiredParamToString(c, param) }
func (h *helper) ParamStr(c echo.Context, param string) (string, error)    { return h.ParamToString(c, param) }
func (h *helper) ReqParamUint(c echo.Context, param string) (uint, error)  { return h.RequiredParamToUint(c, param) }
func (h *helper) ParamUint(c echo.Context, param string) (uint, error)     { return h.ParamToUint(c, param) }
func (h *helper) ReqQueryStr(c echo.Context, param string) (string, error) { return h.RequiredQueryToString(c, param) }
func (h *helper) QueryStr(c echo.Context, param string) (string, error)    { return h.QueryToString(c, param) }
func (h *helper) ReqQueryUint(c echo.Context, param string) (uint, error)  { return h.RequiredQueryToUint(c, param) }
func (h *helper) QueryUint(c echo.Context, param string) (uint, error)     { return h.QueryToUint(c, param) }
func (h *helper) QueryBool(c echo.Context, param string) (bool, error)     { return h.QueryToBool(c, param) }
func (h *helper) FormToBool(c echo.Context, param string) (bool, error)    { return h.FormBool(c, param) }

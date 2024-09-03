package helperecho

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *helperEcho) Bind(c echo.Context, i interface{}) error {
	if err := c.Bind(&i); err != nil {
		return err
	}

	return nil
}

func (h *helperEcho) RequiredFormToString(c echo.Context, param string) (string, error) {
	return h.formToString(c, param, true)
}

func (h *helperEcho) FormToString(c echo.Context, param string) (string, error) {
	return h.formToString(c, param, false)
}

func (h *helperEcho) formToString(c echo.Context, param string, required bool) (string, error) {
	str := c.FormValue(param)
	if str == "" && required {
		return "", fmt.Errorf("form field '%s' is required", param)
	}

	return str, nil
}

func (h *helperEcho) RequiredFormToUint(c echo.Context, param string) (uint, error) {
	return h.formToUint(c, param, true)
}

func (h *helperEcho) FormToUint(c echo.Context, param string) (uint, error) {
	return h.formToUint(c, param, false)
}

func (h *helperEcho) formToUint(c echo.Context, param string, required bool) (uint, error) {
	strID := c.FormValue(param)
	if strID == "" {
		if required {
			return 0, fmt.Errorf("form parameter '%s' is required", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helperEcho) ParamToString(c echo.Context, param string) (string, error) {
	return h.paramToString(c, param, false)
}

func (h *helperEcho) RequiredParamToString(c echo.Context, param string) (string, error) {
	return h.paramToString(c, param, true)
}

func (h *helperEcho) paramToString(c echo.Context, param string, required bool) (string, error) {
	str := c.Param(param)
	if str == "" && required {
		return "", fmt.Errorf("path parameter '%s' is required", param)
	}

	return str, nil
}

func (h *helperEcho) RequiredParamToUint(c echo.Context, param string) (uint, error) {
	return h.paramToUint(c, param, true)
}
func (h *helperEcho) ParamToUint(c echo.Context, param string) (uint, error) {
	return h.paramToUint(c, param, false)
}

func (h *helperEcho) paramToUint(c echo.Context, param string, required bool) (uint, error) {
	strID := c.Param(param)
	if strID == "" || strID == "0" {
		if required {
			return 0, fmt.Errorf("path parameter '%s' is required", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helperEcho) RequiredQueryToString(c echo.Context, param string) (string, error) {
	return h.queryToString(c, param, true)
}

func (h *helperEcho) QueryToString(c echo.Context, param string) (string, error) {
	return h.queryToString(c, param, false)
}

func (h *helperEcho) queryToString(c echo.Context, param string, required bool) (string, error) {
	str := strings.TrimSpace(c.QueryParam(param))
	if str == "" && required {
		return "", fmt.Errorf("query parameter '%s' is required", param)
	}

	return str, nil
}

func (h *helperEcho) RequiredQueryToUint(c echo.Context, param string) (uint, error) {
	return h.queryToUint(c, param, true)
}

func (h *helperEcho) QueryToUint(c echo.Context, param string) (uint, error) {
	return h.queryToUint(c, param, false)
}

func (h *helperEcho) queryToUint(c echo.Context, param string, required bool) (uint, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return 0, fmt.Errorf("query parameter '%s' is required", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helperEcho) QueryToBool(c echo.Context, param string) (bool, error) {
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
		return false, fmt.Errorf("query parameter '%s' has an invalid value, need to be 'true', '1', 'false', '0' or ''", param)
	}
}

func (h *helperEcho) convStrToUint(str string, required bool) (uint, error) {
	ID, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		if required {
			return 0, err
		}

		return 0, err
	}

	return uint(ID), nil
}

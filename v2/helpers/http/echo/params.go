package helperecho

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Valdenirmezadri/core-go/v2/context/corectx"
	"github.com/labstack/echo/v4"
)

func (h *helper) Bind(c echo.Context, i interface{}) error {
	if err := c.Bind(&i); err != nil {
		return err
	}

	return nil
}

func (h *helper) RequiredFormToString(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.formToString(ctx, c, param, true)
}

func (h *helper) FormToString(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.formToString(ctx, c, param, false)
}

func (h *helper) formToString(ctx corectx.UserContext, c echo.Context, param string, required bool) (string, error) {
	str := c.FormValue(param)
	if str == "" && required {
		return "", h.errT(ctx, "FormRequired", param)
	}

	return str, nil
}

func (h *helper) RequiredFormToUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.formToUint(ctx, c, param, true)
}

func (h *helper) FormToUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.formToUint(ctx, c, param, false)
}

func (h *helper) formToUint(ctx corectx.UserContext, c echo.Context, param string, required bool) (uint, error) {
	strID := c.FormValue(param)
	if strID == "" {
		if required {
			return 0, h.errT(ctx, "FormRequired", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helper) RequiredFormToUintArray(ctx corectx.UserContext, c echo.Context, param string) ([]uint, error) {
	return h.formToUintArray(ctx, c, param, true)
}

func (h *helper) FormToUintArray(ctx corectx.UserContext, c echo.Context, param string) ([]uint, error) {
	return h.formToUintArray(ctx, c, param, false)
}

func (h *helper) formToUintArray(ctx corectx.UserContext, c echo.Context, param string, required bool) ([]uint, error) {
	values, err := c.FormParams()
	if err != nil {
		return nil, err
	}

	strs := values[param]
	if len(strs) == 0 {
		if required {
			return nil, h.errT(ctx, "FormRequired", param)
		}

		return []uint{}, nil
	}

	ids := make([]uint, 0, len(strs))
	for _, str := range strs {
		id, err := h.convStrToUint(str, true)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (h *helper) ParamToString(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.paramToString(ctx, c, param, false)
}

func (h *helper) RequiredParamToString(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.paramToString(ctx, c, param, true)
}

func (h *helper) paramToString(ctx corectx.UserContext, c echo.Context, param string, required bool) (string, error) {
	str := c.Param(param)
	if str == "" && required {
		return "", h.errT(ctx, "PathRequired", param)
	}

	return str, nil
}

func (h *helper) RequiredParamToUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.paramToUint(ctx, c, param, true)
}
func (h *helper) ParamToUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.paramToUint(ctx, c, param, false)
}

func (h *helper) paramToUint(ctx corectx.UserContext, c echo.Context, param string, required bool) (uint, error) {
	strID := c.Param(param)
	if strID == "" || strID == "0" {
		if required {
			return 0, h.errT(ctx, "PathRequired", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helper) RequiredQueryToString(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.queryToString(ctx, c, param, true)
}

func (h *helper) QueryToString(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.queryToString(ctx, c, param, false)
}

func (h *helper) queryToString(ctx corectx.UserContext, c echo.Context, param string, required bool) (string, error) {
	str := strings.TrimSpace(c.QueryParam(param))
	if str == "" && required {
		return "", h.errT(ctx, "QueryRequired", param)
	}

	return str, nil
}

func (h *helper) RequiredQueryToUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.queryToUint(ctx, c, param, true)
}

func (h *helper) QueryToUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.queryToUint(ctx, c, param, false)
}

func (h *helper) queryToUint(ctx corectx.UserContext, c echo.Context, param string, required bool) (uint, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return 0, h.errT(ctx, "QueryRequired", param)
		}

		return 0, nil
	}

	return h.convStrToUint(strID, required)
}

func (h *helper) RequiredQueryToUint16(ctx corectx.UserContext, c echo.Context, param string) (uint16, error) {
	return h.queryToUint16(ctx, c, param, true)
}

func (h *helper) QueryToUint16(ctx corectx.UserContext, c echo.Context, param string) (uint16, error) {
	return h.queryToUint16(ctx, c, param, false)
}

func (h *helper) queryToUint16(ctx corectx.UserContext, c echo.Context, param string, required bool) (uint16, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return 0, h.errT(ctx, "QueryRequired", param)
		}

		return 0, nil
	}

	return h.convStrToUint16(strID)
}

func (h *helper) FormBool(ctx corectx.UserContext, c echo.Context, param string) (bool, error) {
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
		return false, h.errT(ctx, "InvalidBoolForm", param)
	}
}

func (h *helper) QueryToBool(ctx corectx.UserContext, c echo.Context, param string) (bool, error) {
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
		return false, h.errT(ctx, "InvalidBoolQuery", param)
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

func (h *helper) convStrToUint16(str string) (uint16, error) {
	ID, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(ID), nil
}

func (h *helper) errT(ctx corectx.UserContext, id, field string) error {
	return errors.New(h.t.S(ctx, id, map[string]string{"Field": field}))
}

// Aliases
func (h *helper) ReqFormStr(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.RequiredFormToString(ctx, c, param)
}
func (h *helper) FormStr(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.FormToString(ctx, c, param)
}
func (h *helper) ReqFormUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.RequiredFormToUint(ctx, c, param)
}
func (h *helper) FormUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.FormToUint(ctx, c, param)
}
func (h *helper) ReqFormUintArr(ctx corectx.UserContext, c echo.Context, param string) ([]uint, error) {
	return h.RequiredFormToUintArray(ctx, c, param)
}
func (h *helper) FormUintArr(ctx corectx.UserContext, c echo.Context, param string) ([]uint, error) {
	return h.FormToUintArray(ctx, c, param)
}
func (h *helper) ReqParamStr(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.RequiredParamToString(ctx, c, param)
}
func (h *helper) ParamStr(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.ParamToString(ctx, c, param)
}
func (h *helper) ReqParamUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.RequiredParamToUint(ctx, c, param)
}
func (h *helper) ParamUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.ParamToUint(ctx, c, param)
}
func (h *helper) ReqQueryStr(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.RequiredQueryToString(ctx, c, param)
}
func (h *helper) QueryStr(ctx corectx.UserContext, c echo.Context, param string) (string, error) {
	return h.QueryToString(ctx, c, param)
}
func (h *helper) ReqQueryUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.RequiredQueryToUint(ctx, c, param)
}
func (h *helper) QueryUint(ctx corectx.UserContext, c echo.Context, param string) (uint, error) {
	return h.QueryToUint(ctx, c, param)
}
func (h *helper) ReqQueryUint16(ctx corectx.UserContext, c echo.Context, param string) (uint16, error) {
	return h.RequiredQueryToUint16(ctx, c, param)
}
func (h *helper) QueryUint16(ctx corectx.UserContext, c echo.Context, param string) (uint16, error) {
	return h.QueryToUint16(ctx, c, param)
}
func (h *helper) QueryBool(ctx corectx.UserContext, c echo.Context, param string) (bool, error) {
	return h.QueryToBool(ctx, c, param)
}
func (h *helper) FormToBool(ctx corectx.UserContext, c echo.Context, param string) (bool, error) {
	return h.FormBool(ctx, c, param)
}

package helperecho

import (
	"encoding/json"
	"errors"
	"fmt"
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

// RequiredParamToStrArr lê o parâmetro repetido (ids=1&ids=2) ou um array JSON
// (ids=["1","2"]) do form ou da query.
// Exige ao menos um valor; valores em branco são ignorados.
func (h *helper) RequiredParamToStrArr(ctx corectx.UserContext, c echo.Context, param string) ([]string, error) {
	return h.paramToStrArr(ctx, c, param, true)
}

// ParamToStrArr é como RequiredParamToStrArr, mas sem valor devolve [] vazio.
func (h *helper) ParamToStrArr(ctx corectx.UserContext, c echo.Context, param string) ([]string, error) {
	return h.paramToStrArr(ctx, c, param, false)
}

// jsonArrToStrs converte um array JSON de textos ou números ("[1,2]") em strings.
func jsonArrToStrs(value string) ([]string, error) {
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.UseNumber()

	items := []any{}
	if err := decoder.Decode(&items); err != nil {
		return nil, err
	}

	strs := make([]string, len(items))
	for i, item := range items {
		switch item.(type) {
		case string, json.Number:
			strs[i] = fmt.Sprint(item)
		default:
			return nil, fmt.Errorf("invalid item %v", item)
		}
	}

	return strs, nil
}

func (h *helper) paramToStrArr(ctx corectx.UserContext, c echo.Context, param string, required bool) ([]string, error) {
	params, err := c.FormParams()
	if err != nil {
		return nil, err
	}

	values := []string{}
	for _, value := range params[param] {
		items := []string{value}
		if strings.HasPrefix(strings.TrimSpace(value), "[") {
			if items, err = jsonArrToStrs(value); err != nil {
				return nil, fmt.Errorf("%s: %w", param, err)
			}
		}

		for _, item := range items {
			if item = strings.TrimSpace(item); item != "" {
				values = append(values, item)
			}
		}
	}

	if len(values) == 0 && required {
		return nil, h.errT(ctx, "ParamRequired", param)
	}

	return values, nil
}

// RequiredParamToArrStrArr lê do form ou da query um parâmetro com array de arrays
// em JSON (groups=[["a","b"],["c","d"]]). Exige ao menos um array.
func (h *helper) RequiredParamToArrStrArr(ctx corectx.UserContext, c echo.Context, param string) ([][]string, error) {
	return paramToArrArr[string](h, ctx, c, param, true)
}

// ParamToArrStrArr é como RequiredParamToArrStrArr, mas sem valor devolve [] vazio.
func (h *helper) ParamToArrStrArr(ctx corectx.UserContext, c echo.Context, param string) ([][]string, error) {
	return paramToArrArr[string](h, ctx, c, param, false)
}

// RequiredParamToArrUintArr é como RequiredParamToArrStrArr, para números
// (ids=[[1,2],[3,4]]). Negativo ou texto dá erro.
func (h *helper) RequiredParamToArrUintArr(ctx corectx.UserContext, c echo.Context, param string) ([][]uint, error) {
	return paramToArrArr[uint](h, ctx, c, param, true)
}

// ParamToArrUintArr é como RequiredParamToArrUintArr, mas sem valor devolve [] vazio.
func (h *helper) ParamToArrUintArr(ctx corectx.UserContext, c echo.Context, param string) ([][]uint, error) {
	return paramToArrArr[uint](h, ctx, c, param, false)
}

func paramToArrArr[T any](h *helper, ctx corectx.UserContext, c echo.Context, param string, required bool) ([][]T, error) {
	values := [][]T{}
	if str := strings.TrimSpace(c.FormValue(param)); str != "" {
		if err := json.Unmarshal([]byte(str), &values); err != nil {
			return nil, fmt.Errorf("%s: %w", param, err)
		}
	}

	if len(values) == 0 && required {
		return nil, h.errT(ctx, "ParamRequired", param)
	}

	return values, nil
}

// RequiredParamToUintArr lê o parâmetro repetido (ids=1&ids=2) ou um array JSON
// (ids=[1,2]) do form ou da query.
// Exige ao menos um valor; em branco e 0 são ignorados, como no ParamToUint.
func (h *helper) RequiredParamToUintArr(ctx corectx.UserContext, c echo.Context, param string) ([]uint, error) {
	return h.paramToUintArr(ctx, c, param, true)
}

// ParamToUintArr é como RequiredParamToUintArr, mas sem valor devolve [] vazio.
func (h *helper) ParamToUintArr(ctx corectx.UserContext, c echo.Context, param string) ([]uint, error) {
	return h.paramToUintArr(ctx, c, param, false)
}

func (h *helper) paramToUintArr(ctx corectx.UserContext, c echo.Context, param string, required bool) ([]uint, error) {
	strs, err := h.paramToStrArr(ctx, c, param, false)
	if err != nil {
		return nil, err
	}

	values := []uint{}
	for _, str := range strs {
		if str == "0" {
			continue
		}

		value, err := h.convStrToUint(str, required)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", param, err)
		}

		values = append(values, value)
	}

	if len(values) == 0 && required {
		return nil, h.errT(ctx, "ParamRequired", param)
	}

	return values, nil
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
func (h *helper) ReqParamStrArr(ctx corectx.UserContext, c echo.Context, param string) ([]string, error) {
	return h.RequiredParamToStrArr(ctx, c, param)
}
func (h *helper) ParamStrArr(ctx corectx.UserContext, c echo.Context, param string) ([]string, error) {
	return h.ParamToStrArr(ctx, c, param)
}
func (h *helper) ReqParamArrStrArr(ctx corectx.UserContext, c echo.Context, param string) ([][]string, error) {
	return h.RequiredParamToArrStrArr(ctx, c, param)
}
func (h *helper) ParamArrStrArr(ctx corectx.UserContext, c echo.Context, param string) ([][]string, error) {
	return h.ParamToArrStrArr(ctx, c, param)
}
func (h *helper) ReqParamArrUintArr(ctx corectx.UserContext, c echo.Context, param string) ([][]uint, error) {
	return h.RequiredParamToArrUintArr(ctx, c, param)
}
func (h *helper) ParamArrUintArr(ctx corectx.UserContext, c echo.Context, param string) ([][]uint, error) {
	return h.ParamToArrUintArr(ctx, c, param)
}
func (h *helper) ReqParamUintArr(ctx corectx.UserContext, c echo.Context, param string) ([]uint, error) {
	return h.RequiredParamToUintArr(ctx, c, param)
}
func (h *helper) ParamUintArr(ctx corectx.UserContext, c echo.Context, param string) ([]uint, error) {
	return h.ParamToUintArr(ctx, c, param)
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

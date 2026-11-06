package helperecho

import (
	"encoding/json"
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

func (h *helperEcho) FormToBool(c echo.Context, param string) (bool, error) {
	return h.FormBool(c, param)
}

func (h *helperEcho) FormBool(c echo.Context, param string) (bool, error) {
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
		return false, fmt.Errorf("form boolean parameter '%s' is required", param)
	}
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

// RequiredParamToStrArr lê o parâmetro repetido (ids=1&ids=2) ou um array JSON
// (ids=["1","2"]) do form ou da query.
// Exige ao menos um valor; valores em branco são ignorados.
func (h *helperEcho) RequiredParamToStrArr(c echo.Context, param string) ([]string, error) {
	return h.paramToStrArr(c, param, true)
}

// ParamToStrArr é como RequiredParamToStrArr, mas sem valor devolve [] vazio.
func (h *helperEcho) ParamToStrArr(c echo.Context, param string) ([]string, error) {
	return h.paramToStrArr(c, param, false)
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

func (h *helperEcho) paramToStrArr(c echo.Context, param string, required bool) ([]string, error) {
	params, err := c.FormParams()
	if err != nil {
		return nil, err
	}

	values := []string{}
	for _, value := range params[param] {
		items := []string{value}
		if strings.HasPrefix(strings.TrimSpace(value), "[") {
			if items, err = jsonArrToStrs(value); err != nil {
				return nil, fmt.Errorf("parameter '%s': %w", param, err)
			}
		}

		for _, item := range items {
			if item = strings.TrimSpace(item); item != "" {
				values = append(values, item)
			}
		}
	}

	if len(values) == 0 && required {
		return nil, fmt.Errorf("parameter '%s' is required", param)
	}

	return values, nil
}

// RequiredParamToArrStrArr lê do form ou da query um parâmetro com array de arrays
// em JSON (groups=[["a","b"],["c","d"]]). Exige ao menos um array.
func (h *helperEcho) RequiredParamToArrStrArr(c echo.Context, param string) ([][]string, error) {
	return paramToArrArr[string](c, param, true)
}

// ParamToArrStrArr é como RequiredParamToArrStrArr, mas sem valor devolve [] vazio.
func (h *helperEcho) ParamToArrStrArr(c echo.Context, param string) ([][]string, error) {
	return paramToArrArr[string](c, param, false)
}

// RequiredParamToArrUintArr é como RequiredParamToArrStrArr, para números
// (ids=[[1,2],[3,4]]). Negativo ou texto dá erro.
func (h *helperEcho) RequiredParamToArrUintArr(c echo.Context, param string) ([][]uint, error) {
	return paramToArrArr[uint](c, param, true)
}

// ParamToArrUintArr é como RequiredParamToArrUintArr, mas sem valor devolve [] vazio.
func (h *helperEcho) ParamToArrUintArr(c echo.Context, param string) ([][]uint, error) {
	return paramToArrArr[uint](c, param, false)
}

func paramToArrArr[T any](c echo.Context, param string, required bool) ([][]T, error) {
	values := [][]T{}
	if str := strings.TrimSpace(c.FormValue(param)); str != "" {
		if err := json.Unmarshal([]byte(str), &values); err != nil {
			return nil, fmt.Errorf("parameter '%s': %w", param, err)
		}
	}

	if len(values) == 0 && required {
		return nil, fmt.Errorf("parameter '%s' is required", param)
	}

	return values, nil
}

// RequiredParamToUintArr lê o parâmetro repetido (ids=1&ids=2) ou um array JSON
// (ids=[1,2]) do form ou da query.
// Exige ao menos um valor; em branco e 0 são ignorados, como no ParamToUint.
func (h *helperEcho) RequiredParamToUintArr(c echo.Context, param string) ([]uint, error) {
	return h.paramToUintArr(c, param, true)
}

// ParamToUintArr é como RequiredParamToUintArr, mas sem valor devolve [] vazio.
func (h *helperEcho) ParamToUintArr(c echo.Context, param string) ([]uint, error) {
	return h.paramToUintArr(c, param, false)
}

func (h *helperEcho) paramToUintArr(c echo.Context, param string, required bool) ([]uint, error) {
	strs, err := h.paramToStrArr(c, param, false)
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
			return nil, fmt.Errorf("parameter '%s': %w", param, err)
		}

		values = append(values, value)
	}

	if len(values) == 0 && required {
		return nil, fmt.Errorf("parameter '%s' is required", param)
	}

	return values, nil
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

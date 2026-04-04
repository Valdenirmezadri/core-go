package helperecho

import (
	"time"

	"github.com/labstack/echo/v4"
)

const (
	layout = "2006-01-02" // Para datas no formato "YYYY-MM-DD"
)

func (h *helper) RequiredQueryToDate(c echo.Context, param string) (time.Time, error) {
	return h._queryToDate(c, param, true)
}

func (h *helper) QueryToDate(c echo.Context, param string) (time.Time, error) {
	return h._queryToDate(c, param, false)
}

func (h *helper) _queryToDate(c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, h.errT("QueryDateRequired", param)
		}

		return time.Time{}, nil
	}

	date, err := h._convStrToDate(strID, param, required)
	if err != nil {
		return time.Time{}, err
	}

	return date.Truncate(24 * time.Hour), nil
}

func (h *helper) _convStrToDate(str, param string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layout, str)
	if err != nil {
		if required {
			return time.Time{}, h.errT("QueryDateInvalid", param)
		}

		return time.Time{}, nil
	}

	return parsedTime, nil
}

// Aliases
func (h *helper) ReqQueryDate(c echo.Context, param string) (time.Time, error) { return h.RequiredQueryToDate(c, param) }
func (h *helper) QueryDate(c echo.Context, param string) (time.Time, error)    { return h.QueryToDate(c, param) }

package helperecho

import (
	"time"

	"github.com/labstack/echo/v4"
)

const (
	layoutWithTime = "2006-01-02T15:04:05" // Para datas no formato "YYYY-MM-DDTHH:MM:SS"
)

func (h *helper) RequiredQueryToDateTime(c echo.Context, param string) (time.Time, error) {
	return h._queryToDateTime(c, param, true)
}

func (h *helper) QueryToDateTime(c echo.Context, param string) (time.Time, error) {
	return h._queryToDateTime(c, param, false)
}

func (h *helper) _queryToDateTime(c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, h.errT("QueryDateTimeRequired", param)
		}

		return time.Time{}, nil
	}

	return h._convStrToDateTime(strID, param, required)
}

func (h *helper) _convStrToDateTime(str, param string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layoutWithTime, str)
	if err != nil {
		if required {
			return time.Time{}, h.errT("QueryDateTimeInvalid", param)
		}

		return time.Time{}, nil
	}

	return parsedTime, nil
}

// Aliases
func (h *helper) ReqQueryDateTime(c echo.Context, param string) (time.Time, error) { return h.RequiredQueryToDateTime(c, param) }
func (h *helper) QueryDateTime(c echo.Context, param string) (time.Time, error)    { return h.QueryToDateTime(c, param) }

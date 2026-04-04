package helperecho

import (
	"time"

	"github.com/labstack/echo/v4"
)

const (
	layoutTime = "15:04:05" // Para tempos no formato "HH:MM:SS"
)

func (h *helper) RequiredQueryToTime(c echo.Context, param string) (time.Time, error) {
	return h._queryToTime(c, param, true)
}

func (h *helper) QueryToTime(c echo.Context, param string) (time.Time, error) {
	return h._queryToTime(c, param, false)
}

func (h *helper) _queryToTime(c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, h.errT("QueryTimeRequired", param)
		}

		return time.Time{}, nil
	}

	return h._convStrToTime(strID, param, required)
}

func (h *helper) _convStrToTime(str, param string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layoutTime, str)
	if err != nil {
		if required {
			return time.Time{}, h.errT("QueryTimeInvalid", param)
		}
		return time.Time{}, nil
	}
	return parsedTime, nil
}

// Aliases
func (h *helper) ReqQueryTime(c echo.Context, param string) (time.Time, error) { return h.RequiredQueryToTime(c, param) }
func (h *helper) QueryTime(c echo.Context, param string) (time.Time, error)    { return h.QueryToTime(c, param) }

package helperecho

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	layout = "2006-01-02" // Para datas no formato "YYYY-MM-DD"
)

func (h *helperEcho) RequiredQueryToDate(c echo.Context, param string) (time.Time, error) {
	return h._queryToDate(c, param, true)
}

func (h *helperEcho) QueryToDate(c echo.Context, param string) (time.Time, error) {
	return h._queryToDate(c, param, false)
}

func (h *helperEcho) _queryToDate(c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, fmt.Errorf("query parameter '%s' is required as 'YYYY-MM-DD'", param)
		}

		return time.Time{}, nil
	}

	date, err := h._convStrToDate(strID, required)
	if err != nil {
		return time.Time{}, err
	}

	return date.Truncate(24 * time.Hour), nil
}

func (h *helperEcho) _convStrToDate(str string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layout, str)
	if err != nil {
		if required {
			return time.Time{}, fmt.Errorf("query parameter '%s' is not a valid date format 'YYYY-MM-DD'", str)
		}

		return time.Time{}, nil
	}

	return parsedTime, nil
}

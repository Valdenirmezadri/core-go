package helperecho

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	layoutWithTime = "2006-01-02T15:04:05" // Para datas no formato "YYYY-MM-DDTHH:MM:SS"
)

func (h *helperEcho) RequiredQueryToDateTime(c echo.Context, param string) (time.Time, error) {
	return h._queryToDateTime(c, param, true)
}

func (h *helperEcho) QueryToDateTime(c echo.Context, param string) (time.Time, error) {
	return h._queryToDateTime(c, param, false)
}

func (h *helperEcho) _queryToDateTime(c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, fmt.Errorf("query parameter '%s' is required as 'YYYY-MM-DDTHH:MM:SS'", param)
		}

		return time.Time{}, nil
	}

	return h._convStrToDateTime(strID, required)
}

func (h *helperEcho) _convStrToDateTime(str string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layoutWithTime, str)
	if err != nil {
		if required {
			return time.Time{}, fmt.Errorf("query parameter '%s' is not a valid date format 'YYYY-MM-DDTHH:MM:SS'", str)
		}

		return time.Time{}, nil
	}

	return parsedTime, nil
}

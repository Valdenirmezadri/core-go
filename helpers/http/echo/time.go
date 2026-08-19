package helperecho

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	layoutTime = "15:04:05" // Para tempos no formato "HH:MM:SS"

)

func (h *helperEcho) RequiredQueryToTime(c echo.Context, param string) (time.Time, error) {
	return h._queryToTime(c, param, true)
}

func (h *helperEcho) QueryToTime(c echo.Context, param string) (time.Time, error) {
	return h._queryToTime(c, param, false)
}

func (h *helperEcho) _queryToTime(c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, fmt.Errorf("query parameter '%s' is required as 'HH:MM:SS'", param)
		}

		return time.Time{}, nil
	}

	return h._convStrToTime(strID, required)
}

func (h *helperEcho) _convStrToTime(str string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layoutTime, str)
	if err != nil {
		if required {
			return time.Time{}, fmt.Errorf("query parameter '%s' time is not a valid format 'HH:MM:SS'", str)
		}

		return time.Time{}, nil
	}
	return parsedTime, nil
}

package helperecho

import (
	"time"

	"github.com/Valdenirmezadri/core-go/v2/context/corectx"
	"github.com/labstack/echo/v4"
)

const (
	layoutTime = "15:04:05" // Para tempos no formato "HH:MM:SS"
)

func (h *helper) RequiredQueryToTime(ctx corectx.UserContext, c echo.Context, param string) (time.Time, error) {
	return h._queryToTime(ctx, c, param, true)
}

func (h *helper) QueryToTime(ctx corectx.UserContext, c echo.Context, param string) (time.Time, error) {
	return h._queryToTime(ctx, c, param, false)
}

func (h *helper) _queryToTime(ctx corectx.UserContext, c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, h.errT(ctx, "QueryTimeRequired", param)
		}

		return time.Time{}, nil
	}

	return h._convStrToTime(ctx, strID, param, required)
}

func (h *helper) _convStrToTime(ctx corectx.UserContext, str, param string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layoutTime, str)
	if err != nil {
		if required {
			return time.Time{}, h.errT(ctx, "QueryTimeInvalid", param)
		}
		return time.Time{}, nil
	}
	return parsedTime, nil
}

// Aliases
func (h *helper) ReqQueryTime(ctx corectx.UserContext, c echo.Context, param string) (time.Time, error) {
	return h.RequiredQueryToTime(ctx, c, param)
}
func (h *helper) QueryTime(ctx corectx.UserContext, c echo.Context, param string) (time.Time, error) {
	return h.QueryToTime(ctx, c, param)
}

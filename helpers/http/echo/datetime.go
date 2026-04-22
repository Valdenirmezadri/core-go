package helperecho

import (
	"time"

	"github.com/Valdenirmezadri/core-go/context/corectx"
	"github.com/labstack/echo/v4"
)

const (
	layoutWithTime = "2006-01-02T15:04:05" // Para datas no formato "YYYY-MM-DDTHH:MM:SS"
)

func (h *helper) RequiredQueryToDateTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error) {
	return h._queryToDateTime(ctx, c, param, true)
}

func (h *helper) QueryToDateTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error) {
	return h._queryToDateTime(ctx, c, param, false)
}

func (h *helper) _queryToDateTime(ctx corectx.Context, c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, h.errT(ctx, "QueryDateTimeRequired", param)
		}

		return time.Time{}, nil
	}

	return h._convStrToDateTime(ctx, strID, param, required)
}

func (h *helper) _convStrToDateTime(ctx corectx.Context, str, param string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layoutWithTime, str)
	if err != nil {
		if required {
			return time.Time{}, h.errT(ctx, "QueryDateTimeInvalid", param)
		}

		return time.Time{}, nil
	}

	return parsedTime, nil
}

// Aliases
func (h *helper) ReqQueryDateTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error) { return h.RequiredQueryToDateTime(ctx, c, param) }
func (h *helper) QueryDateTime(ctx corectx.Context, c echo.Context, param string) (time.Time, error)    { return h.QueryToDateTime(ctx, c, param) }

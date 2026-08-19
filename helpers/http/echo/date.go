package helperecho

import (
	"time"

	"github.com/Valdenirmezadri/core-go/context/corectx"
	"github.com/labstack/echo/v4"
)

const (
	layout = "2006-01-02" // Para datas no formato "YYYY-MM-DD"
)

func (h *helper) RequiredQueryToDate(ctx corectx.UserContext, c echo.Context, param string) (time.Time, error) {
	return h._queryToDate(ctx, c, param, true)
}

func (h *helper) QueryToDate(ctx corectx.UserContext, c echo.Context, param string) (time.Time, error) {
	return h._queryToDate(ctx, c, param, false)
}

func (h *helper) _queryToDate(ctx corectx.UserContext, c echo.Context, param string, required bool) (time.Time, error) {
	strID := c.QueryParam(param)
	if strID == "" {
		if required {
			return time.Time{}, h.errT(ctx, "QueryDateRequired", param)
		}

		return time.Time{}, nil
	}

	date, err := h._convStrToDate(ctx, strID, param, required)
	if err != nil {
		return time.Time{}, err
	}

	return date.Truncate(24 * time.Hour), nil
}

func (h *helper) _convStrToDate(ctx corectx.UserContext, str, param string, required bool) (time.Time, error) {
	parsedTime, err := time.Parse(layout, str)
	if err != nil {
		if required {
			return time.Time{}, h.errT(ctx, "QueryDateInvalid", param)
		}

		return time.Time{}, nil
	}

	return parsedTime, nil
}

// Aliases
func (h *helper) ReqQueryDate(ctx corectx.UserContext, c echo.Context, param string) (time.Time, error) {
	return h.RequiredQueryToDate(ctx, c, param)
}
func (h *helper) QueryDate(ctx corectx.UserContext, c echo.Context, param string) (time.Time, error) {
	return h.QueryToDate(ctx, c, param)
}

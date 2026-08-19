package helperecho

import (
	"errors"

	"github.com/Valdenirmezadri/core-go/v2/context/corectx"
	"github.com/labstack/echo/v4"
)

func (h *helper) RequiredPropToString(ctx corectx.UserContext, c echo.Context, propName string) (string, error) {
	return h.PropToString(ctx, c, propName, true)
}

func (h *helper) PropToString(ctx corectx.UserContext, c echo.Context, propName string, required bool) (string, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return "", err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return "", h.errT(ctx, "PropRequired", propName)
		}

		return "", nil
	}

	str, ok := prop.(string)
	if !ok {
		return "", h.propTypeErr(ctx, propName, "string")
	}

	if str == "" && required {
		return "", h.errT(ctx, "PropRequired", propName)
	}

	return str, nil
}

func (h *helper) RequiredPropToUint(ctx corectx.UserContext, c echo.Context, propName string) (uint, error) {
	return h.PropToUint(ctx, c, propName, true)
}

func (h *helper) PropToUint(ctx corectx.UserContext, c echo.Context, propName string, required bool) (uint, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return 0, err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return 0, h.errT(ctx, "PropRequired", propName)
		}
		return 0, nil
	}

	f, ok := prop.(float64)
	if !ok {
		return 0, h.propTypeErr(ctx, propName, "number")
	}

	if f == 0 && required {
		return 0, h.errT(ctx, "PropRequired", propName)
	}

	return uint(f), nil
}

func (h *helper) RequiredPropToBool(ctx corectx.UserContext, c echo.Context, propName string) (bool, error) {
	return h.PropToBool(ctx, c, propName, true)
}

func (h *helper) PropToBool(ctx corectx.UserContext, c echo.Context, propName string, required bool) (bool, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return false, err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return false, h.errT(ctx, "PropNotFound", propName)
		}

		return false, nil
	}

	b, ok := prop.(bool)
	if !ok {
		return false, h.propTypeErr(ctx, propName, "boolean")
	}

	return b, nil
}

func (h *helper) propTypeErr(ctx corectx.UserContext, propName, kind string) error {
	return errors.New(h.t.S(ctx, "PropType", map[string]string{
		"Field": propName,
		"Type":  kind,
	}))
}

// Aliases
func (h *helper) ReqPropStr(ctx corectx.UserContext, c echo.Context, propName string) (string, error) {
	return h.RequiredPropToString(ctx, c, propName)
}
func (h *helper) PropStr(ctx corectx.UserContext, c echo.Context, propName string, required bool) (string, error) {
	return h.PropToString(ctx, c, propName, required)
}
func (h *helper) ReqPropUint(ctx corectx.UserContext, c echo.Context, propName string) (uint, error) {
	return h.RequiredPropToUint(ctx, c, propName)
}
func (h *helper) PropUint(ctx corectx.UserContext, c echo.Context, propName string, required bool) (uint, error) {
	return h.PropToUint(ctx, c, propName, required)
}
func (h *helper) ReqPropBool(ctx corectx.UserContext, c echo.Context, propName string) (bool, error) {
	return h.RequiredPropToBool(ctx, c, propName)
}
func (h *helper) PropBool(ctx corectx.UserContext, c echo.Context, propName string, required bool) (bool, error) {
	return h.PropToBool(ctx, c, propName, required)
}

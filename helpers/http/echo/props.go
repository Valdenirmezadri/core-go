package helperecho

import (
	"errors"

	"github.com/labstack/echo/v4"
)

func (h *helper) RequiredPropToString(c echo.Context, propName string) (string, error) {
	return h.PropToString(c, propName, true)
}

func (h *helper) PropToString(c echo.Context, propName string, required bool) (string, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return "", err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return "", h.errT("PropRequired", propName)
		}

		return "", nil
	}

	str, ok := prop.(string)
	if !ok {
		return "", h.propTypeErr(propName, "string")
	}

	if str == "" && required {
		return "", h.errT("PropRequired", propName)
	}

	return str, nil
}

func (h *helper) RequiredPropToUint(c echo.Context, propName string) (uint, error) {
	return h.PropToUint(c, propName, true)
}

func (h *helper) PropToUint(c echo.Context, propName string, required bool) (uint, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return 0, err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return 0, h.errT("PropRequired", propName)
		}
		return 0, nil
	}

	f, ok := prop.(float64)
	if !ok {
		return 0, h.propTypeErr(propName, "number")
	}

	if f == 0 && required {
		return 0, h.errT("PropRequired", propName)
	}

	return uint(f), nil
}

func (h *helper) RequiredPropToBool(c echo.Context, propName string) (bool, error) {
	return h.PropToBool(c, propName, true)
}

func (h *helper) PropToBool(c echo.Context, propName string, required bool) (bool, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return false, err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return false, h.errT("PropNotFound", propName)
		}

		return false, nil
	}

	b, ok := prop.(bool)
	if !ok {
		return false, h.propTypeErr(propName, "boolean")
	}

	return b, nil
}

func (h *helper) propTypeErr(propName, kind string) error {
	return errors.New(h.t.S(nil, "PropType", map[string]string{
		"Field": propName,
		"Type":  kind,
	}))
}

// Aliases
func (h *helper) ReqPropStr(c echo.Context, propName string) (string, error)          { return h.RequiredPropToString(c, propName) }
func (h *helper) PropStr(c echo.Context, propName string, required bool) (string, error) { return h.PropToString(c, propName, required) }
func (h *helper) ReqPropUint(c echo.Context, propName string) (uint, error)            { return h.RequiredPropToUint(c, propName) }
func (h *helper) PropUint(c echo.Context, propName string, required bool) (uint, error) { return h.PropToUint(c, propName, required) }
func (h *helper) ReqPropBool(c echo.Context, propName string) (bool, error)            { return h.RequiredPropToBool(c, propName) }
func (h *helper) PropBool(c echo.Context, propName string, required bool) (bool, error) { return h.PropToBool(c, propName, required) }

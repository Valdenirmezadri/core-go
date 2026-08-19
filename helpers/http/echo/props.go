package helperecho

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (h *helperEcho) RequiredPropToString(c echo.Context, propName string) (string, error) {
	return h.PropToString(c, propName, true)
}

func (h *helperEcho) PropToString(c echo.Context, propName string, required bool) (string, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return "", err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return "", h.propIsRequired(propName)
		}

		return "", nil
	}

	str, ok := prop.(string)
	if !ok {
		return "", h.propIsNot(propName, "string")
	}

	if str == "" && required {
		return "", h.propIsRequired(propName)
	}

	return str, nil
}

func (h *helperEcho) RequiredPropToUint(c echo.Context, propName string) (uint, error) {
	return h.PropToUint(c, propName, true)
}

func (h *helperEcho) PropToUint(c echo.Context, propName string, required bool) (uint, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return 0, err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return 0, h.propIsRequired(propName)
		}
		return 0, nil
	}

	f, ok := prop.(float64)
	if !ok {
		return 0, h.propIsNot(propName, "number")
	}

	if f == 0 && required {
		return 0, h.propIsRequired(propName)
	}

	return uint(f), nil
}

func (h *helperEcho) RequiredPropToBool(c echo.Context, propName string) (bool, error) {
	return h.PropToBool(c, propName, true)
}

func (h *helperEcho) PropToBool(c echo.Context, propName string, required bool) (bool, error) {
	data := make(map[string]any)

	if err := c.Bind(&data); err != nil {
		return false, err
	}

	prop, ok := data[propName]
	if !ok {
		if required {
			return false, h.propNotFound(propName)
		}

		return false, nil
	}

	b, ok := prop.(bool)
	if !ok {
		return false, h.propIsNot(propName, "boolean")
	}

	return b, nil
}

// propIsNot return 'property "propName" is required'
func (h *helperEcho) propIsRequired(propName string) error {
	return fmt.Errorf(`property "%s" is required`, propName)
}

// propIsNot return 'property "propName" notfound'
func (h *helperEcho) propNotFound(propName string) error {
	return fmt.Errorf(`property "%s" not found`, propName)
}

// propIsNot return 'property "propName" need to be "kind"'
func (h *helperEcho) propIsNot(propName, kind string) error {
	return fmt.Errorf(`property "%s" need to be "%s"`, propName, kind)
}

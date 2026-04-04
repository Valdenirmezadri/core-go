package jwtdata

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const ctxUserID = "user_id"
const ctxUserKind = "user_kind"

func Lang(c echo.Context) string {
	language := c.Request().Header.Get("Accept-Language")
	if language == "" {
		language = "pt-BR"
		c.Request().Header.Set("Accept-Language", language)
	}

	return language
}

func GetJWTData(c echo.Context) (uint, uint8, string, error) {
	userID, ok := c.Get(ctxUserID).(uint)
	if !ok {
		return 0, 0, "", fmt.Errorf("user_id required")
	}

	userKind, ok := c.Get(ctxUserKind).(uint8)
	if !ok {
		return 0, 0, "", fmt.Errorf("user_kind required")
	}

	return userID, userKind, Lang(c), nil
}

func SetUserID(c echo.Context, userID uint) {
	c.Set(ctxUserID, userID)
}

func SetUserKind(c echo.Context, kind uint8) {
	c.Set(ctxUserKind, kind)
}

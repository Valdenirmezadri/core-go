package secure

import (
	"github.com/Valdenirmezadri/core-go/hterr"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type SecureHandler interface {
	Lang(c echo.Context) string
	TokenFromHeader(htErr hterr.Err, c echo.Context, key []byte) (token *jwt.Token, err error)
	IsRefresh(htErr hterr.Err, c echo.Context, token *jwt.Token) (bool, error)
	Validate(htErr hterr.Err, c echo.Context, token *jwt.Token) (echo.Context, error)
	ValidateRefresh(htErr hterr.Err, c echo.Context, token *jwt.Token) (echo.Context, error)
	NewTokens(ID uint, kind uint8, isProd bool, key []byte) (accessToken, refreshToken string, err error)
}

type handler struct {
}

func New() SecureHandler {
	return &handler{}
}

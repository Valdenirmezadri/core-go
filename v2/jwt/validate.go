package secure

import (
	"github.com/Valdenirmezadri/core-go/v2/hterr"
	jwtdata "github.com/Valdenirmezadri/core-go/v2/jwt/jwt-data"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *handler) Lang(c echo.Context) string {
	return jwtdata.Lang(c)
}

func (h *handler) IsRefresh(htErr hterr.Err, c echo.Context, token *jwt.Token) (bool, error) {
	_tokenType, err := h.getUint8FromClaims(htErr, token, tokenTypeKey)
	if err != nil {
		return false, err
	}

	tokenType := TokenType(_tokenType)
	if err := tokenType.valid(htErr); err != nil {
		return false, err
	}

	if !tokenType.Refresh() {
		return false, echo.ErrUnauthorized
	}

	return true, nil
}

func (h *handler) TokenFromHeader(htErr hterr.Err, c echo.Context, key []byte) (token *jwt.Token, err error) {
	if c.RealIP() == "" {
		return nil, htErr.JWTInvalidIP(nil)
	}

	auth := c.Request().Header.Get("authorization")
	token, err = h.ValidToken(htErr, auth, key)
	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return nil, echo.ErrUnauthorized
}

func (h *handler) Validate(htErr hterr.Err, c echo.Context, token *jwt.Token) (echo.Context, error) {
	_tokenType, err := h.getUint8FromClaims(htErr, token, tokenTypeKey)
	if err != nil {
		return nil, err
	}

	tokenType := TokenType(_tokenType)
	if err := tokenType.valid(htErr); err != nil {
		return nil, err
	}

	if !tokenType.Access() {
		return nil, echo.ErrUnauthorized
	}

	if err := h.recoverFromToken(htErr, c, token); err != nil {
		return nil, err
	}

	return c, nil
}

func (h *handler) ValidateRefresh(htErr hterr.Err, c echo.Context, token *jwt.Token) (echo.Context, error) {
	_tokenType, err := h.getUint8FromClaims(htErr, token, tokenTypeKey)
	if err != nil {
		return nil, err
	}

	tokenType := TokenType(_tokenType)
	if err := tokenType.valid(htErr); err != nil {
		return nil, err
	}

	if !tokenType.Refresh() {
		return nil, echo.ErrUnauthorized
	}

	if err := h.recoverFromToken(htErr, c, token); err != nil {
		return nil, err
	}

	return c, nil
}

func (h *handler) recoverFromToken(htErr hterr.Err, c echo.Context, token *jwt.Token) error {
	userID, err := h.getUintFromClaims(htErr, token, idKey)
	if err != nil {
		return err
	}

	userType, err := h.getUint8FromClaims(htErr, token, userTypeKey)
	if err != nil {
		return err
	}

	jwtdata.SetUserID(c, userID)
	jwtdata.SetUserKind(c, userType)

	return nil
}

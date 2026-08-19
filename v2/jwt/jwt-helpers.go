package secure

import (
	"github.com/Valdenirmezadri/core-go/hterr"
	"github.com/golang-jwt/jwt/v5"
)

func (h *handler) getClaims(htErr hterr.Err, token *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, htErr.JWTClaimsNotFound(nil)
}

func (h *handler) getUintFromClaims(htErr hterr.Err, token *jwt.Token, param string) (uint, error) {
	claims, err := h.getClaims(htErr, token)
	if err != nil {
		return 0, err
	}

	if rawID, ok := claims[param]; ok {
		if ID, ok := rawID.(float64); ok {
			return uint(ID), nil
		}
	}

	return 0, htErr.JWTKeyNotFound(nil)
}

func (h *handler) getUint8FromClaims(htErr hterr.Err, token *jwt.Token, param string) (uint8, error) {
	n, err := h.getUintFromClaims(htErr, token, param)
	if err != nil {
		return 0, err
	}

	if n <= 255 {
		return uint8(n), nil
	}

	return 0, htErr.JWTClaimsNotFound(nil)
}

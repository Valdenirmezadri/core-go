package secure

import (
	"strings"
	"time"

	"github.com/Valdenirmezadri/core-go/hterr"
	"github.com/golang-jwt/jwt/v5"
)

const (
	idKey        = "u-id"
	userTypeKey  = "u-type"
	tokenTypeKey = "t-type"
	expKey       = "exp"
)

func (h handler) howLong(isProd bool) time.Duration {
	if isProd {
		return time.Hour * 10
	}

	return time.Hour * 24
}

func (h handler) NewTokens(ID uint, kind uint8, isProd bool, key []byte) (accessToken, refreshToken string, err error) {
	accessToken, err = h.newToken(ID, kind, AccessToken, time.Now().Add(h.howLong(isProd)), key)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = h.newToken(ID, kind, RefreshToken, time.Now().Add(time.Hour*168), key)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (h *handler) newToken(ID uint, userKind uint8, tokenType TokenType, expiration time.Time, key []byte) (jwtToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[idKey] = ID
	claims[userTypeKey] = userKind
	claims[tokenTypeKey] = tokenType
	claims[expKey] = expiration.Unix()

	return token.SignedString(key)
}

// ValidToken usado para validar token
func (h *handler) ValidToken(htErr hterr.Err, tokenJWT string, key []byte) (*jwt.Token, error) {
	if len(tokenJWT) == 0 {
		return nil, htErr.JWTInvalid(nil)
	}

	parts := strings.Split(tokenJWT, " ")
	if len(parts) != 2 || parts[1] == "" {
		return nil, htErr.TokenInvalid(nil)
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, htErr.JWTInvalid(nil)
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		if _, ok := token.Claims.(jwt.MapClaims); ok {
			return token, nil
		}
	}

	return nil, htErr.JWTClaimsNotFound(nil)
}

func (h handler) GetUint(htErr hterr.Err, claims jwt.MapClaims, k string) (uint, error) {
	floatClaim, ok := claims[k].(float64)
	if !ok {
		return 0, htErr.JWTClaimsNotFound(nil)
	}

	uintClaim := uint(floatClaim)
	if floatClaim < 0 || float64(uintClaim) != floatClaim {
		return 0, htErr.JWTClaimsNotFound(nil)
	}

	return uintClaim, nil
}

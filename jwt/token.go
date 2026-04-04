package secure

import (
	"fmt"

	"github.com/Valdenirmezadri/core-go/hterr"
)

type TokenType uint8

func (t TokenType) valid(err hterr.Err) error {
	if t.Access() || t.Refresh() {
		return nil
	}

	return err.TokenInvalid(nil)
}

func (t TokenType) String() string {
	if t.Access() {
		return "access"
	}

	if t.Refresh() {
		return "refresh"
	}

	return fmt.Sprintf(`token type "%d" not valid`, t)
}

func (t TokenType) Access() bool {
	return t == AccessToken
}

func (t TokenType) Refresh() bool {
	return t == RefreshToken
}

const (
	AccessToken TokenType = iota + 1
	RefreshToken
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

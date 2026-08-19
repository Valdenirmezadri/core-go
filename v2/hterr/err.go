package hterr

import (
	"embed"
	"errors"
	"fmt"

	"github.com/Valdenirmezadri/core-go/v2/context/corectx"
	"github.com/Valdenirmezadri/core-go/v2/i18n"
	"gorm.io/gorm"
)

//go:embed locales/*.toml
var localeFS embed.FS

// Sentinel errors for errors.Is() comparison
var (
	ErrJWTInvalid        = errors.New("jwt_invalid")
	ErrJWTKeyNotFound    = errors.New("jwt_key_not_found")
	ErrJWTClaimsNotFound = errors.New("jwt_claims_not_found")
	ErrJWTInvalidIP      = errors.New("jwt_invalid_ip")
)

var (
	ErrTokenInvalid       = errors.New("token_invalid")
	ErrTokenBlock         = errors.New("token_block")
	ErrTokenNotAuthorized = errors.New("token_not_authorized")
)

var (
	ErrEmptyIP  = errors.New("empty_ip")
	ErrBannedIP = errors.New("banned_ip")
)

var (
	ErrUserIDRequired   = errors.New("user_id_required")
	ErrUserKindRequired = errors.New("user_kind_required")
)

var (
	ErrDBConnectionNil           = errors.New("db_connection_nil")
	ErrTransactionAlreadyRunning = errors.New("transaction_already_running")
	ErrTransactionEmpty          = errors.New("transaction_empty")
)

var (
	ErrRecipientNotSubscribed = errors.New("recipient_not_subscribed")
)

var (
	ErrEmailInvalid     = errors.New("email_invalid")
	ErrPasswordMismatch = errors.New("password_mismatch")
)

type Err interface {
	JWTInvalid(ctx corectx.Context) error
	JWTKeyNotFound(ctx corectx.Context) error
	JWTClaimsNotFound(ctx corectx.Context) error
	JWTInvalidIP(ctx corectx.Context) error
	TokenInvalid(ctx corectx.Context) error
	TokenBlock(ctx corectx.Context) error
	TokenNotAuthorized(ctx corectx.Context) error
	EmptyIP(ctx corectx.Context) error
	BannedIP(ctx corectx.Context) error
	UserIDRequired(ctx corectx.Context) error
	UserKindRequired(ctx corectx.Context) error
	DBConnectionNil(ctx corectx.Context) error
	TransactionAlreadyRunning(ctx corectx.Context) error
	TransactionEmpty(ctx corectx.Context) error
	NotSubscribed(ctx corectx.Context, ID string) error
	EmailInvalid(ctx corectx.Context) error
	PasswordMismatch(ctx corectx.Context) error
	RecordNotFound(err error) bool
	Is(err error, targets ...error) bool
	As(err error, target any) bool
	EvalAs(err error, target any, fn func() bool) bool
	Errorf(format string, a ...any) error
}

type htErr struct {
	t i18n.Switch
}

func New(logger func() i18n.Logger) (Err, error) {
	t, err := i18n.New(localeFS, logger)
	if err != nil {
		return nil, err
	}

	return &htErr{t: t}, nil
}

func (h *htErr) JWTInvalid(ctx corectx.Context) error {
	return h.wrap(ctx, "JWTInvalid", ErrJWTInvalid)
}

func (h *htErr) JWTKeyNotFound(ctx corectx.Context) error {
	return h.wrap(ctx, "JWTKeyNotFound", ErrJWTKeyNotFound)
}

func (h *htErr) JWTClaimsNotFound(ctx corectx.Context) error {
	return h.wrap(ctx, "JWTClaimsNotFound", ErrJWTClaimsNotFound)
}

func (h *htErr) JWTInvalidIP(ctx corectx.Context) error {
	return h.wrap(ctx, "JWTInvalidIP", ErrJWTInvalidIP)
}

func (h *htErr) TokenInvalid(ctx corectx.Context) error {
	return h.wrap(ctx, "TokenInvalid", ErrTokenInvalid)
}

func (h *htErr) TokenBlock(ctx corectx.Context) error {
	return h.wrap(ctx, "TokenBlock", ErrTokenBlock)
}

func (h *htErr) EmptyIP(ctx corectx.Context) error {
	return h.wrap(ctx, "EmptyIP", ErrEmptyIP)
}

func (h *htErr) BannedIP(ctx corectx.Context) error {
	return h.wrap(ctx, "BannedIP", ErrBannedIP)
}

func (h *htErr) UserIDRequired(ctx corectx.Context) error {
	return h.wrap(ctx, "UserIDRequired", ErrUserIDRequired)
}

func (h *htErr) UserKindRequired(ctx corectx.Context) error {
	return h.wrap(ctx, "UserKindRequired", ErrUserKindRequired)
}

func (h *htErr) TokenNotAuthorized(ctx corectx.Context) error {
	return h.wrap(ctx, "TokenNotAuthorized", ErrTokenNotAuthorized)
}

func (h *htErr) DBConnectionNil(ctx corectx.Context) error {
	return h.wrap(ctx, "DBConnectionNil", ErrDBConnectionNil)
}

func (h *htErr) TransactionAlreadyRunning(ctx corectx.Context) error {
	return h.wrap(ctx, "TransactionAlreadyRunning", ErrTransactionAlreadyRunning)
}

func (h *htErr) TransactionEmpty(ctx corectx.Context) error {
	return h.wrap(ctx, "TransactionEmpty", ErrTransactionEmpty)
}

func (h *htErr) NotSubscribed(ctx corectx.Context, ID string) error {
	msg := h.t.S(ctx, "NotSubscribed", map[string]string{"ID": ID})
	return fmt.Errorf("%s: %w", msg, ErrRecipientNotSubscribed)
}

func (h *htErr) EmailInvalid(ctx corectx.Context) error {
	return h.wrap(ctx, "EmailInvalid", ErrEmailInvalid)
}

func (h *htErr) PasswordMismatch(ctx corectx.Context) error {
	return h.wrap(ctx, "PasswordMismatch", ErrPasswordMismatch)
}

func (h *htErr) wrap(ctx corectx.Context, msgID string, sentinel error) error {
	return fmt.Errorf("%s: %w", h.t.S(ctx, msgID, nil), sentinel)
}

func (h *htErr) RecordNotFound(err error) bool {
	return h.Is(err, gorm.ErrRecordNotFound)
}

func (h *htErr) Is(err error, targets ...error) bool {
	if err == nil {
		return false
	}

	for _, target := range targets {
		if ok := errors.Is(err, target); ok {
			return true
		}
	}

	return false
}

func (h *htErr) As(err error, target any) bool {
	return errors.As(err, target)
}

// EvalAs faz errors.As para target e, se bem-sucedido, chama fn com o valor extraído.
// Retorna false se err não puder ser desempacotado em target.
//
// Exemplo:
//
//	var httpErr whatsmeow.DownloadHTTPError
//	retry := h.err.EvalAs(err, &httpErr, func() bool {
//	    return httpErr.StatusCode == 403 || httpErr.StatusCode == 401
//	})
func (h *htErr) EvalAs(err error, target any, fn func() bool) bool {
	if !h.As(err, target) {
		return false
	}
	return fn()
}

func (h *htErr) Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

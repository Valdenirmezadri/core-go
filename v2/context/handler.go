package appctx

import (
	"context"

	jwtdata "github.com/Valdenirmezadri/core-go/jwt/jwt-data"
	"github.com/Valdenirmezadri/core-go/pgorm"
	"github.com/labstack/echo/v4"
)

type ContextHandler interface {
	SetOnEcho(c echo.Context, ctx ContextHandler)
	GetFromEcho(c echo.Context) ContextHandler
	//NewGuestFromEcho(c echo.Context, userID uint, userKind uint8, lang string) (Context, error)
	NewFromEcho(c echo.Context) (Context, error)
	New(ctx context.Context, lang string, userID uint, kind uint8) (Context, error)
	Background(ctx context.Context, lang string, userID uint, kind uint8) (Context, error)
	NewSuper(userID uint, userKind uint8, lang string) (Context, error)
	Tools() Tools
}

type handler struct {
	conn  pgorm.DB
	tools Tools
}

func (h *handler) Tools() Tools {
	return h.tools
}

func (h *handler) NewSuper(userID uint, userKind uint8, lang string) (Context, error) {
	return newContext(context.Background(), lang, userID, userKind, h.conn, h.tools, false), nil
}

func (h *handler) NewGuestFromEcho(c echo.Context, userID uint, userKind uint8, lang string) (Context, error) {
	return newContext(c.Request().Context(), lang, userID, userKind, h.conn, h.tools, true), nil
}

func (h *handler) NewFromEcho(c echo.Context) (Context, error) {
	userID, userKInd, lang, err := jwtdata.GetJWTData(c)
	if err != nil {
		return nil, err
	}

	return h.New(c.Request().Context(), lang, userID, userKInd)
}

func (h *handler) New(ctx context.Context, lang string, userID uint, kind uint8) (Context, error) {
	return h.new(ctx, lang, userID, kind, true)
}

func (h *handler) Background(ctx context.Context, lang string, userID uint, kind uint8) (Context, error) {
	return h.new(ctx, lang, userID, kind, false)
}

func (h *handler) new(ctx context.Context, lang string, userID uint, kind uint8, useContext bool) (Context, error) {
	if userID == 0 {
		return nil, h.tools.Err().UserIDRequired(nil)
	}

	if kind == 0 {
		return nil, h.tools.Err().UserKindRequired(nil)
	}

	return newContext(ctx, lang, userID, kind, h.conn, h.tools, useContext), nil
}

func New(conn pgorm.DB, tools Tools) (ContextHandler, error) {
	if conn == nil {
		return nil, tools.Err().DBConnectionNil(nil)
	}

	return &handler{
		conn:  conn,
		tools: tools,
	}, nil
}

func (h *handler) GetFromEcho(c echo.Context) ContextHandler {
	return c.Get("coreCtx").(ContextHandler)
}

func (h *handler) SetOnEcho(c echo.Context, ctx ContextHandler) {
	c.Set("coreCtx", ctx)
}

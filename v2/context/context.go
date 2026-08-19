package appctx

import (
	"context"
	"time"

	"github.com/Valdenirmezadri/core-go/context/corectx"
	"github.com/Valdenirmezadri/core-go/i18n"
	"github.com/Valdenirmezadri/core-go/pgorm"
	"github.com/Valdenirmezadri/core-go/safe"
	"gorm.io/gorm"
)

func newUser(lang string, userID uint, kind uint8) corectx.User {
	language := i18n.Language(lang)
	if err := language.Valid(); err != nil {
		language = i18n.PtBR
	}

	return corectx.User{
		Lang:      language.String(),
		ID:        userID,
		Kind:      kind,
		StartedAt: time.Now(),
	}
}

type Context interface {
	corectx.Context
	Tools() Tools
}

type transaction struct {
	id   string
	gorm *gorm.DB
}

type htContext struct {
	user       corectx.User
	tools      Tools
	useContext bool
	ctx        context.Context
	conn       pgorm.DB
	tx         safe.Item[*transaction]
}

func newContext(ctx context.Context, lang string, userID uint, kind uint8, conn pgorm.DB, tools Tools, useContext bool) *htContext {
	return &htContext{
		user:       newUser(lang, userID, kind),
		tools:      tools,
		useContext: useContext,
		ctx:        ctx,
		conn:       conn,
		tx:         safe.NewItem[*transaction](),
	}
}

func (c *htContext) User() corectx.User { return c.user }
func (c *htContext) Tools() Tools       { return c.tools }

func (c *htContext) Context() context.Context {
	if c.useContext {
		return c.ctx
	}

	return context.Background()
}

func (c *htContext) db() *gorm.DB {
	ctx := c.Context()
	tx := c.tx.Get()
	if tx != nil {
		return tx.gorm.WithContext(ctx)
	}

	return c.conn.Conn().WithContext(ctx)
}

func (c *htContext) Read() *gorm.DB { return c.db() }

func (c *htContext) Write() *gorm.DB { return c.db() }

func (c *htContext) Fetch(opts ...pgorm.QueryOption) *gorm.DB {
	return c.conn.Fetch(c.Context(), opts...)
}

func (c *htContext) BeginTransaction(id string) {
	conn := c.conn.Conn()
	if conn == nil {
		return
	}

	if c.tx.Get() != nil {
		return
	}

	gorm := conn.WithContext(c.ctx).Begin()

	c.tx.Set(&transaction{id, gorm})
}

func (c *htContext) RollbackTransaction(id string) {
	tx := c.tx.Get()
	if tx == nil || tx.id != id {
		return
	}

	tx.gorm.Rollback()

	c.tx.Set(nil)
}

func (c *htContext) RollbackIfErr(id string, err error) {
	r := recover()
	if err != nil || r != nil {
		c.RollbackTransaction(id)
	}

}

func (c *htContext) CommitTransaction(id string) {
	tx := c.tx.Get()
	if tx == nil || tx.id != id {
		return
	}

	tx.gorm.Commit()

	c.tx.Set(nil)
}

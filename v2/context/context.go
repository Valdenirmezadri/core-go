package appctx

import (
	"context"
	"time"

	"github.com/Valdenirmezadri/core-go/v2/context/corectx"
	"github.com/Valdenirmezadri/core-go/v2/context/transaction"
	"github.com/Valdenirmezadri/core-go/v2/i18n"
	"github.com/Valdenirmezadri/core-go/v2/pgorm"
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

type htContext struct {
	user       corectx.User
	tools      Tools
	useContext bool
	ctx        context.Context
	conn       pgorm.DB
	tx         transaction.Transaction
}

func newContext(ctx context.Context, lang string, userID uint, kind uint8, conn pgorm.DB, tools Tools, useContext bool) *htContext {
	return &htContext{
		user:       newUser(lang, userID, kind),
		tools:      tools,
		useContext: useContext,
		ctx:        ctx,
		conn:       conn,
		// lazy: o Context nasce antes de a conexão existir, e pedi-la aqui daria
		// nil para sempre
		tx: transaction.NewLazy(conn.Conn),
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

func (c *htContext) db() *gorm.DB { return c.tx.DB(c.Context()) }

func (c *htContext) Read() *gorm.DB { return c.db() }

func (c *htContext) Write() *gorm.DB { return c.db() }

func (c *htContext) Fetch(opts ...pgorm.QueryOption) *gorm.DB {
	return c.conn.Fetch(c.Context(), opts...)
}

/*
O controle da transação vive todo no pacote transaction: aqui é só repasse, para
o Context continuar oferecendo as mesmas quatro operações de sempre.
*/
func (c *htContext) BeginTransaction(id string) {
	c.tx.BeginTransaction(c.Context(), id)
}

func (c *htContext) RollbackTransaction(id string) {
	c.tx.RollbackTransaction(id)
}

/*
RollbackIfErr repassa, e o recover fica aqui porque só alcança o panic quando é a
própria função adiada que o chama: recuperar lá dentro do pacote devolveria nulo.

O erro vai por ponteiro para o defer poder ser direto:

	defer ctx.RollbackIfErr(id, &err)

Por valor exigiria embrulhar num closure para o erro ser lido no return em vez da
hora do defer — e o closure afasta o recover um frame, que é o bastante para o
panic passar batido e a transação vazar aberta.
*/
func (c *htContext) RollbackIfErr(id string, err *error) {
	c.tx.RollbackOnErr(id, err, recover())
}

func (c *htContext) RollbackOnErr(id string, err *error, recovered any) {
	c.tx.RollbackOnErr(id, err, recovered)
}

func (c *htContext) CommitTransaction(id string) error {
	return c.tx.CommitTransaction(id)
}

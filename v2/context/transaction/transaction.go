/*
Package transaction guarda e controla a transação do gorm.

É a única implementação de Begin/Rollback/Commit do core: o appctx delega para
cá, e quem não tem um Context do core na mão — repositório que recebe
context.Context da stdlib e a conexão pelo construtor — monta o seu com New.

A transação fica guardada aqui, e DB devolve ela: quem consulta o banco não
precisa saber que está dentro de uma. O id é de quem abriu, e só ele confirma ou
desfaz — assim um trecho que abre transação por conta própria, chamado de dentro
de outra, não corta a de fora pela metade.

	t := transaction.New(conn)

	t.BeginTransaction(ctx, "criar-empresa")
	defer t.RollbackIfErr("criar-empresa", &err)

	// ... passos, cada um consultando por t.DB(ctx) ...

	return t.CommitTransaction("criar-empresa")

Cada Transaction guarda uma transação por vez, então precisa de um por linha de
execução — o appctx cria o dele por request. Um só, guardado num struct de vida
longa e atendido por várias goroutines, misturaria a transação de um request com
as consultas de outro.
*/
package transaction

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/Valdenirmezadri/core-go/v2/htl"
	"github.com/Valdenirmezadri/core-go/v2/safe"
	"gorm.io/gorm"
)

type Transaction interface {
	/*BeginTransaction abre a transação em nome de id.

	Com uma já aberta, não abre outra: o trecho de dentro participa da de fora, e
	só quem abriu confirma. Sem conexão também não abre — as consultas seguem
	direto no banco, que é o que acontecia antes de alguém pedir transação.
	*/
	BeginTransaction(ctx context.Context, id string)

	// RollbackTransaction desfaz a transação, se id for quem a abriu
	RollbackTransaction(id string)

	/*RollbackIfErr desfaz a transação quando o erro apontado por err não é nulo,
	ou quando houve panic. É o que dispensa cada return de lembrar do rollback.

	O erro vai por ponteiro, e não por valor, porque `defer f(id, err)` avalia os
	argumentos na hora do defer — quando err ainda é nulo. Por valor, a função
	nunca enxergaria o erro que o return deixou, e a transação ficaria aberta.

	O panic vira erro de retorno em vez de subir: se subisse, a função sairia com
	o retorno zerado e quem chamou receberia erro nulo como se tivesse gravado.
	Exige, para isso, que o retorno de erro seja nomeado.

	Precisa ser chamada direto pelo defer, senão o recover não alcança o panic:

		defer t.RollbackIfErr(id, &err)
	*/
	RollbackIfErr(id string, err *error)

	/*RollbackOnErr é o RollbackIfErr para quem já chamou recover.

	Existe porque recover só alcança o panic quando é chamado pela própria função
	adiada. Quem oferece o RollbackIfErr por uma camada acima — o appctx faz isso
	— recupera lá e manda o resultado para cá, e assim a decisão continua num
	lugar só.
	*/
	RollbackOnErr(id string, err *error, recovered any)

	/*CommitTransaction confirma a transação, se id for quem a abriu.

	Devolve erro porque commit que falha não gravou nada: engolir aqui deixaria a
	função seguir e devolver sucesso com o banco intacto. Dentro de uma aninhada
	devolve nulo sem confirmar — quem confirma é quem abriu.
	*/
	CommitTransaction(id string) error

	/*DB é por onde as consultas passam.

	Devolve a transação quando há uma aberta, e a conexão quando não há — então a
	mesma consulta serve nos dois casos, sem quem consulta saber em qual está.
	*/
	DB(ctx context.Context) *gorm.DB
}

// New monta o controle sobre uma conexão já resolvida
func New(conn *gorm.DB) Transaction {
	return NewLazy(func() *gorm.DB { return conn })
}

/*
NewLazy monta o controle sobre uma conexão resolvida na hora do uso.

É o que o appctx precisa: lá o Context nasce antes de a conexão existir, e pedi-la
na montagem daria nil para sempre.
*/
func NewLazy(conn func() *gorm.DB) Transaction {
	return &transaction{conn: conn, tx: safe.NewItem[*openTx]()}
}

type openTx struct {
	id string
	db *gorm.DB
}

type transaction struct {
	conn func() *gorm.DB
	tx   safe.Item[*openTx]
}

func (t *transaction) BeginTransaction(ctx context.Context, id string) {
	if t.tx.Get() != nil {
		return
	}

	conn := t.conn()
	if conn == nil {
		return
	}

	tx := conn.WithContext(ctx).Begin()
	if tx.Error != nil {
		return
	}

	t.tx.Set(&openTx{id: id, db: tx})
}

func (t *transaction) RollbackTransaction(id string) {
	tx := t.tx.Get()
	if tx == nil || tx.id != id {
		return
	}

	tx.db.Rollback()

	t.tx.Set(nil)
}

func (t *transaction) RollbackIfErr(id string, err *error) {
	t.RollbackOnErr(id, err, recover())
}

func (t *transaction) RollbackOnErr(id string, err *error, recovered any) {
	if recovered == nil && (err == nil || *err == nil) {
		return
	}

	t.RollbackTransaction(id)

	if recovered == nil {
		return
	}

	// sem ponteiro onde escrever, o panic volta a subir: engolir devolveria
	// sucesso falso para quem chamou
	if err == nil {
		panic(recovered)
	}

	logPanic(id, recovered)

	*err = fmt.Errorf("%v", recovered)
}

/*
logPanic manda a stack para o log.

A stack sai no log, e não no erro: virando erro, o panic perde o rastro que o
recover do servidor logaria, e quem chama recebe "invalid memory address" sem uma
linha de onde. O erro é para o cliente, a stack é para quem vai consertar.

O recover próprio é porque htl.Log() estoura quando ninguém iniciou o log — e
aqui não é lugar de estourar, já estamos tratando um panic. Log que não sai é
menos grave do que transação que não desfaz.
*/
func logPanic(id string, recovered any) {
	defer func() { _ = recover() }()

	htl.Log().Errorf("transação %s desfeita por panic: %v\n%s", id, recovered, debug.Stack())
}

func (t *transaction) CommitTransaction(id string) error {
	tx := t.tx.Get()
	if tx == nil || tx.id != id {
		return nil
	}

	// limpa antes de conferir o erro: commit que falha já encerrou a transação no
	// servidor, e guardá-la faria o rollback seguinte tentar desfazer o que não
	// existe mais
	t.tx.Set(nil)

	return tx.db.Commit().Error
}

func (t *transaction) DB(ctx context.Context) *gorm.DB {
	if tx := t.tx.Get(); tx != nil {
		return tx.db.WithContext(ctx)
	}

	conn := t.conn()
	if conn == nil {
		return nil
	}

	return conn.WithContext(ctx)
}

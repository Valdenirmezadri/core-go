package transaction

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const id = "teste"

var falhaNoPasso = errors.New("passo falhou")

/*
banco devolve um gorm ligado a um banco de mentira.

sqlmock, e não um banco de verdade, porque o que estes testes conferem é
exatamente qual comando saiu — BEGIN, COMMIT ou ROLLBACK. Um banco real
responderia certo aos três e não diria qual aconteceu.
*/
func banco(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("não montou o banco de mentira: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	conn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
		// sem isto o gorm pergunta a versão do servidor ao abrir, e o mock não
		// espera essa consulta
		DriverName:           "postgres",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("não abriu o gorm: %v", err)
	}

	return conn, mock
}

// gravou é um passo qualquer dentro da transação, para haver o que desfazer
func gravou(tx *gorm.DB) error {
	return tx.Exec("UPDATE contas SET nome = ?", "x").Error
}

func esperaGravacao(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta("UPDATE contas")).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func conferir(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("o que saiu para o banco não bate: %v", err)
	}
}

// -------------------- o caminho comum --------------------

// sem erro nenhum, confirma
func TestCommit(t *testing.T) {
	conn, mock := banco(t)

	mock.ExpectBegin()
	esperaGravacao(mock)
	mock.ExpectCommit()

	gravar := func() (err error) {
		tr := New(conn)

		tr.BeginTransaction(context.Background(), id)
		defer tr.RollbackIfErr(id, &err)

		if err := gravou(tr.DB(context.Background())); err != nil {
			return err
		}

		return tr.CommitTransaction(id)
	}

	if err := gravar(); err != nil {
		t.Fatalf("gravação sem problema devolveu erro: %v", err)
	}

	conferir(t, mock)
}

/*
Passo devolvendo erro desfaz, e é o defer que desfaz — nenhum return lembra do
rollback.

O ponteiro é o que faz isto funcionar: o defer avalia os argumentos na hora em
que é escrito, quando err ainda é nulo. Por valor, RollbackIfErr receberia nulo
e a transação ficaria aberta.
*/
func TestRollbackQuandoOPassoFalha(t *testing.T) {
	conn, mock := banco(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	gravar := func() (err error) {
		tr := New(conn)

		tr.BeginTransaction(context.Background(), id)
		defer tr.RollbackIfErr(id, &err)

		return falhaNoPasso
	}

	if err := gravar(); !errors.Is(err, falhaNoPasso) {
		t.Fatalf("esperava o erro do passo, veio: %v", err)
	}

	conferir(t, mock)
}

// -------------------- RollbackIfErr: adiado direto --------------------

/*
TestRollbackIfErr_PanicViraErro é o caso de quem pode ser adiado direto.

A função adiada é o próprio RollbackIfErr, então o recover dele está no frame
certo: pega o panic, desfaz, escreve no erro nomeado, e a função retorna normal
em vez de o panic subir.
*/
func TestRollbackIfErr_PanicViraErro(t *testing.T) {
	conn, mock := banco(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	gravar := func() (err error) {
		tr := New(conn)

		tr.BeginTransaction(context.Background(), id)
		defer tr.RollbackIfErr(id, &err)

		panic("estourou no meio")
	}

	err := gravar()
	if err == nil {
		t.Fatal("o panic saiu como sucesso: quem chamou receberia erro nulo")
	}

	if err.Error() != "estourou no meio" {
		t.Errorf("o erro não diz o que aconteceu: %v", err)
	}

	conferir(t, mock)
}

// -------------------- RollbackOnErr: com camada no meio --------------------

/*
camada é uma casca por cima do Transaction, como o appctx e o Session do propaga.

O RollbackIfErr dela é o que vai no defer de quem usa, então é aqui que o recover
precisa acontecer — o Transaction está um frame abaixo, e lá ele voltaria nulo.
Por isso a camada recupera e entrega o resultado pronto ao RollbackOnErr.
*/
type camada struct{ tr Transaction }

func (c camada) RollbackIfErr(id string, err *error) {
	c.tr.RollbackOnErr(id, err, recover())
}

/*
camadaIngenua é a mesma casca escrita do jeito errado: repassa o RollbackIfErr em
vez de recuperar.

Existe para o teste provar o estrago, e não só descrevê-lo.
*/
type camadaIngenua struct{ tr Transaction }

func (c camadaIngenua) RollbackIfErr(id string, err *error) {
	c.tr.RollbackIfErr(id, err)
}

/*
TestRollbackOnErr_CamadaRecuperaEDesfaz: com o recover no frame certo, a camada
se comporta igualzinho ao caso sem camada.
*/
func TestRollbackOnErr_CamadaRecuperaEDesfaz(t *testing.T) {
	conn, mock := banco(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	gravar := func() (err error) {
		c := camada{tr: New(conn)}

		c.tr.BeginTransaction(context.Background(), id)
		defer c.RollbackIfErr(id, &err)

		panic("estourou no meio")
	}

	err := gravar()
	if err == nil {
		t.Fatal("o panic saiu como sucesso mesmo com a camada recuperando")
	}

	conferir(t, mock)
}

/*
TestCamadaSemRecover_PanicEscapaEDeixaATransacaoAberta é o motivo de RollbackOnErr
existir.

recover só devolve o panic quando é a própria função adiada que o chama. Com a
camada repassando, ele fica dois frames fundo demais e volta nulo: o panic sobe
sem ninguém desfazer nada, e a transação continua aberta segurando conexão e
travas até o banco desistir.

Compila, roda, passa em revisão. Só aparece em produção.
*/
func TestCamadaSemRecover_PanicEscapaEDeixaATransacaoAberta(t *testing.T) {
	conn, mock := banco(t)

	// nenhum ExpectRollback: o teste afirma que o rollback NÃO acontece
	mock.ExpectBegin()

	gravar := func() (err error) {
		c := camadaIngenua{tr: New(conn)}

		c.tr.BeginTransaction(context.Background(), id)
		defer c.RollbackIfErr(id, &err)

		panic("estourou no meio")
	}

	escapou := func() (p any) {
		defer func() { p = recover() }()

		_ = gravar()

		return nil
	}()

	if escapou == nil {
		t.Fatal("o panic foi engolido pela camada, e nem virou erro nem subiu")
	}

	conferir(t, mock)
}

// -------------------- dono da transação --------------------

/*
Quem abriu é quem confirma.

Um trecho que abre transação por conta própria, chamado de dentro de outra, não
pode confirmar a de fora — os passos que ainda vinham ficariam do lado de fora do
commit. Aqui o de dentro tenta confirmar e não acontece nada; quem fecha é o de
fora.
*/
func TestSoQuemAbriuConfirma(t *testing.T) {
	conn, mock := banco(t)

	mock.ExpectBegin()
	esperaGravacao(mock)
	mock.ExpectCommit()

	tr := New(conn)
	ctx := context.Background()

	tr.BeginTransaction(ctx, "de-fora")

	// o de dentro não abre outra, e o commit dele não vale
	tr.BeginTransaction(ctx, "de-dentro")
	tr.CommitTransaction("de-dentro")

	// a gravação vem depois de propósito: se o commit de dentro tivesse valido,
	// ela cairia fora da transação — que é justamente o estrago que a regra evita
	if err := gravou(tr.DB(ctx)); err != nil {
		t.Fatalf("o commit de dentro fechou a transação de fora: %v", err)
	}

	if err := tr.CommitTransaction("de-fora"); err != nil {
		t.Fatalf("o commit de quem abriu falhou: %v", err)
	}

	conferir(t, mock)
}

/*
Só quem abriu desfaz.

Mesma razão do commit, do outro lado: um trecho de dentro desfazendo levaria
junto tudo o que o de fora já tinha feito, e o de fora seguiria achando que sua
transação continua de pé.
*/
func TestSoQuemAbriuDesfaz(t *testing.T) {
	conn, mock := banco(t)

	mock.ExpectBegin()
	mock.ExpectCommit()

	tr := New(conn)
	ctx := context.Background()

	tr.BeginTransaction(ctx, "de-fora")

	// não é o dono: não pode desfazer
	tr.RollbackTransaction("de-dentro")

	tr.CommitTransaction("de-fora")

	conferir(t, mock)
}

/*
Commit que falha vira erro, e não sucesso.

É o silêncio mais caro dos três: erro no passo pelo menos aparece, panic pelo
menos estoura, mas commit engolido deixa a função devolver sucesso com o banco
exatamente como estava antes.
*/
func TestCommitQueFalhaViraErro(t *testing.T) {
	conn, mock := banco(t)

	falhaNoCommit := errors.New("conexão caiu no commit")

	mock.ExpectBegin()
	esperaGravacao(mock)
	mock.ExpectCommit().WillReturnError(falhaNoCommit)

	gravar := func() (err error) {
		tr := New(conn)

		tr.BeginTransaction(context.Background(), id)
		defer tr.RollbackIfErr(id, &err)

		if err := gravou(tr.DB(context.Background())); err != nil {
			return err
		}

		return tr.CommitTransaction(id)
	}

	if err := gravar(); !errors.Is(err, falhaNoCommit) {
		t.Fatalf("o commit falhou e a função devolveu: %v", err)
	}

	conferir(t, mock)
}

// Sem conexão não abre transação, e DB não devolve nada para consultar — em vez
// de estourar na primeira consulta
func TestSemConexao(t *testing.T) {
	tr := New(nil)

	tr.BeginTransaction(context.Background(), id)

	if db := tr.DB(context.Background()); db != nil {
		t.Error("devolveu algo para consultar sem ter conexão")
	}
}

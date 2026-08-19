// Package withcontext fornece utilitários para executar operações
// com verificação automática de cancelamento de contexto.
//
// As funções deste pacote verificam se o contexto foi cancelado
// antes e/ou durante a execução, retornando ctx.Err() quando apropriado.
//
// Exemplo de uso:
//
//	var lastID uint
//	err := withcontext.Paginate(
//		ctx,
//		func() ([]Item, error) {
//			return repo.Fetch(lastID, 100)
//		},
//		func(item Item) error {
//			lastID = item.ID
//			return process(item)
//		},
//	)
package withcontext

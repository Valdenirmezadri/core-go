package corectx

import (
	"context"
	"time"

	"github.com/Valdenirmezadri/core-go/v2/pgorm"
	"gorm.io/gorm"
)

type User struct {
	Lang      string
	ID        uint
	Kind      uint8
	StartedAt time.Time
}

type UserContext interface {
	User() User
}

type Context interface {
	UserContext
	Context() context.Context
	Write() *gorm.DB
	Read() *gorm.DB
	Fetch(...pgorm.QueryOption) *gorm.DB
	BeginTransaction(id string)
	RollbackTransaction(id string)
	RollbackIfErr(id string, err *error)
	// RollbackOnErr é o RollbackIfErr para quem já chamou recover — é o que uma
	// camada que embrulha o Context precisa, porque recover só alcança o panic
	// quando é a própria função adiada que o chama
	RollbackOnErr(id string, err *error, recovered any)
	CommitTransaction(id string) error
}

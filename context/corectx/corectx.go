package corectx

import (
	"context"
	"time"

	"github.com/Valdenirmezadri/core-go/pgorm"
	"gorm.io/gorm"
)

type User struct {
	Lang      string
	ID        uint
	Kind      uint8
	StartedAt time.Time
}

type Context interface {
	User() User
	Context() context.Context
	Write() *gorm.DB
	Read() *gorm.DB
	Fetch(...pgorm.QueryOption) *gorm.DB
	BeginTransaction(id string)
	RollbackTransaction(id string)
	RollbackIfErr(id string, err error)
	CommitTransaction(id string)
}

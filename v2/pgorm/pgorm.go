package pgorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Valdenirmezadri/core-go/v2/environment"
	htl "github.com/Valdenirmezadri/core-go/v2/htl"
	"github.com/Valdenirmezadri/core-go/v2/safe"
	"github.com/hashicorp/go-version"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB interface {
	Conn() *gorm.DB
	Read(ctx context.Context) *gorm.DB
	Write(ctx context.Context) *gorm.DB
	Fetch(ctx context.Context, opts ...QueryOption) *gorm.DB
	Version() (version *version.Version, err error)
	NewVersion(ver string) (version *version.Version, err error)
	UpdateVersion(newVer *version.Version) error
	initDB() error
	LogLevel(env environment.Environment, l string)
}

type conn struct {
	_gorm safe.Item[*gorm.DB]
}

func (c *conn) Conn() *gorm.DB {
	return c._gorm.Get()
}

func (c *conn) Read(ctx context.Context) *gorm.DB {
	return c.Conn().WithContext(ctx)
}

func (c *conn) Fetch(ctx context.Context, opts ...QueryOption) *gorm.DB {
	return applyOptions(c.Read(ctx), opts)
}

func (c *conn) Write(ctx context.Context) *gorm.DB {
	return c.Conn().WithContext(ctx)
}

func new(g *gorm.DB) DB {
	return &conn{_gorm: safe.NewItemWithData(g)}
}

type DBConfig struct {
	Host         string
	SSL          string
	TimeZone     string `mapstructure:"time_zone"`
	OpenConns    uint   `mapstructure:"open_conns"`
	MaxIdleConns uint   `mapstructure:"maxIdle_conns"`
	Name         string
	User         string
	Pass         string
	Log          string
}

func (c *conn) LogLevel(env environment.Environment, l string) {
	lvl := level(l)

	c._gorm.Update(func(g *gorm.DB) *gorm.DB {
		config := *g.Config
		config.Logger = newLogger(env, lvl)
		g.Config = &config
		return g
	})
}

func level(s string) logger.LogLevel {
	switch strings.ToUpper(s) {
	case "SILENT":
		return logger.Silent
	case "ERROR":
		return logger.Error
	case "WARN":
		return logger.Warn
	case "INFO":
		return logger.Info
	default:
		return logger.Info
	}
}

func New(host, sslMode, timeZone string, openConns, maxIdleConns uint, dbname, user, pass string, logLevel string, env environment.Environment) (DB, error) {
	return connectGORM(host, user, pass, dbname, sslMode, timeZone, openConns, maxIdleConns, level(logLevel), env)
}

func buildDBURI(HOST, USER, PASS, DBNAME, sslMode, timeZone string) string {
	return fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s TimeZone=%s",
		HOST, USER, DBNAME, PASS, sslMode, timeZone)
}

// ConnectGORM Abre a conexão com o banco de dados
func connectGORM(HOST, USER, PASS, DBNAME, sslMode, timeZone string, openConns, maxIdleConns uint, logLevel logger.LogLevel, env environment.Environment) (DB, error) {
	htl.Log().Debugf("conectando ao banco %+v\n", DBNAME)

	db, err := gorm.Open(postgres.Open(buildDBURI(HOST, USER, PASS, DBNAME, sslMode, timeZone)), &gorm.Config{
		FullSaveAssociations: false,
		AllowGlobalUpdate:    false,
		Logger:               newLogger(env, logLevel),
		CreateBatchSize:      100,
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(int(openConns))
	sqlDB.SetMaxIdleConns(int(maxIdleConns))
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	conn := new(db)

	if err := conn.initDB(); err != nil {
		return nil, err
	}

	return conn, nil
}

func newLogger(env environment.Environment, loglevel logger.LogLevel) logger.Interface {
	color := env == environment.Developer
	ignoreNotfound := true
	if loglevel == logger.Info {
		ignoreNotfound = false
	}

	return logger.New(htl.Log(), logger.Config{
		SlowThreshold:             1 * time.Second,
		LogLevel:                  loglevel,
		IgnoreRecordNotFoundError: ignoreNotfound,
		Colorful:                  color,
	})

}

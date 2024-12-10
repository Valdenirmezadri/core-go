package pgorm

import (
	"fmt"
	"strings"
	"sync"
	"time"

	htl "github.com/Valdenirmezadri/core-go/htl"
	"github.com/Valdenirmezadri/htrelay-server/modules/core"
	"github.com/hashicorp/go-version"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB interface {
	Conn() *gorm.DB
	Version() (version *version.Version, err error)
	NewVersion(ver string) (version *version.Version, err error)
	UpdateVersion(newVer *version.Version) error
	initDB() error
	LogLevel(env core.Environment, l string)
}

type conn struct {
	lock *sync.RWMutex
	gorm *gorm.DB
}

func (c conn) Conn() *gorm.DB {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.gorm
}

func new(g *gorm.DB) DB {
	return &conn{
		lock: &sync.RWMutex{},
		gorm: g,
	}
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

func (c *conn) LogLevel(env core.Environment, l string) {
	level := level(l)

	c.lock.Lock()
	config := *c.gorm.Config
	config.Logger = newLogger(env, level)
	c.gorm.Config = &config
	c.lock.Unlock()
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

func New(host, sslMode, timeZone string, openConns, maxIdleConns uint, dbname, user, pass string, logLevel string, env core.Environment) (DB, error) {
	return connectGORM(host, user, pass, dbname, sslMode, timeZone, openConns, maxIdleConns, level(logLevel), env)
}

func buildDBURI(HOST, USER, PASS, DBNAME, sslMode, timeZone string) string {
	return fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s TimeZone=%s",
		HOST, USER, DBNAME, PASS, sslMode, timeZone)
}

// ConnectGORM Abre a conexão com o banco de dados
func connectGORM(HOST, USER, PASS, DBNAME, sslMode, timeZone string, openConns, maxIdleConns uint, logLevel logger.LogLevel, env core.Environment) (DB, error) {
	htl.Log().Debugf("conectando ao banco %+v\n", DBNAME)

	db, err := gorm.Open(postgres.Open(buildDBURI(HOST, USER, PASS, DBNAME, sslMode, timeZone)), &gorm.Config{
		FullSaveAssociations: false,
		AllowGlobalUpdate:    false,
		Logger:               newLogger(env, logLevel),
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

func newLogger(env core.Environment, loglevel logger.LogLevel) logger.Interface {
	color := env == core.Developer
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

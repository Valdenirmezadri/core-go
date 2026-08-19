package pgorm

import (
	"fmt"

	htl "github.com/Valdenirmezadri/core-go/htl"
	version "github.com/hashicorp/go-version"
	"github.com/jackc/pgx/v5/pgconn"
)

type database struct {
	ID      uint8  `gorm:"primaryKey"`
	Version string `json:"version"`
}

func (c conn) NewVersion(ver string) (*version.Version, error) {
	return version.NewVersion(ver)
}

func (c conn) Version() (*version.Version, error) {
	data := database{ID: 1}
	if err := c.Conn().First(&data).Error; err != nil {
		return nil, err
	}

	return c.NewVersion(data.Version)
}

func (c conn) UpdateVersion(new *version.Version) error {
	data := database{ID: 1}
	if err := c.Conn().First(&data).Error; err != nil {
		return err
	}

	current, err := version.NewVersion(data.Version)
	if err != nil {
		return err
	}

	if current.Equal(new) {
		return fmt.Errorf("new version %s is equal as current %s", new, current)
	}

	return c.Conn().Model(&database{ID: 1}).Update("version", new.String()).Error
}

func (c conn) initDB() error {
	err := c.Conn().First(&database{}).Error
	if err == nil {
		return nil
	}

	pgErr := err.(*pgconn.PgError)
	if pgErr.Code == "42P01" {
		htl.Log().Info("iniciando banco de dados...")
		if err := c.unaccent(); err != nil {
			return err
		}

		if err := c.dataBase(); err != nil {
			return err
		}

		data := database{ID: 1, Version: "0.0.1"}
		if err := c.Conn().Create(&data).Error; err != nil {
			return err
		}

		initVer, err := c.NewVersion("0.1")
		if err != nil {
			return err
		}

		return c.UpdateVersion(initVer)
	}

	return err
}

func (c conn) unaccent() error {
	return c.Conn().Exec("CREATE EXTENSION IF NOT EXISTS unaccent").Error
}

func (c conn) dataBase() error {
	return c.Conn().Debug().AutoMigrate(&database{})
}

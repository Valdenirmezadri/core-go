package pgorm

import "gorm.io/gorm"

type QueryOption func(*gorm.DB) *gorm.DB

func Order(asc bool) QueryOption {
	if asc {
		return func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }
	}
	return func(db *gorm.DB) *gorm.DB { return db.Order("id DESC") }
}

func Limit(n int) QueryOption {
	if n != 0 {
		return func(db *gorm.DB) *gorm.DB { return db.Limit(n) }
	}

	return func(db *gorm.DB) *gorm.DB { return db.Limit(100) }
}

func applyOptions(db *gorm.DB, opts []QueryOption) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

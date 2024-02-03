package repo

import (
	"digital-book/internal/types"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Migrate() error {
	for _, table := range []any{types.User{}, types.Mark{}, types.Subject{}} {
		if err := r.db.AutoMigrate(table); err != nil {
			return err
		}
	}
	return nil
}

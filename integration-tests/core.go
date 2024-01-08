package integrationtests_test

import (
	"context"
	"digital-book/internal/config"
	"digital-book/internal/repo"
	"digital-book/internal/service"
	"digital-book/internal/types"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type srv interface {
	CreateSubject(ctx context.Context, subject *types.Subject) error
	GetSubject(ctx context.Context, id string) (res types.Subject, err error)

	GetMarks(ctx context.Context) (res []types.MarkResponse, err error)
	DeleteMark(ctx context.Context, id string) error
}

func runInTransaction(out func(s srv)) {
	db, err := gorm.Open(postgres.Open(config.Config.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	db.Transaction(func(tx *gorm.DB) error {
		rep := repo.New(tx)
		out(service.New(rep))
		return errors.New("rollback")
	})
}

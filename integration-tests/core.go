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
	CreateStudent(ctx context.Context, student *types.Student) error
	GetStudents(ctx context.Context) (res []types.Student, err error)
	GetStudent(ctx context.Context, id string) (res types.Student, err error)
	DeleteStudent(ctx context.Context, id string) error
	UpdateStudent(ctx context.Context, name, class, id string) error

	CreateSubject(ctx context.Context, subject *types.Subject) error
	GetSubject(ctx context.Context, id string) (res types.Subject, err error)
	GetSubjects(ctx context.Context) (res []types.Subject, err error)
	DeleteSubject(ctx context.Context, id string) error
	UpdateSubject(ctx context.Context, title, id string) error

	GetMarks(ctx context.Context) (res []types.MarkResponse, err error)
	DeleteMark(ctx context.Context, id string) error
	CreateMark(ctx context.Context, mark *types.Mark) error
	GetMark(ctx context.Context, id string) (res types.Mark, err error)
	UpdateMark(ctx context.Context, studentID, subjectID, value, id string) error
	AvgMarks(ctx context.Context) (res []types.MarkAverege, err error)
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

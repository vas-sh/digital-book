package service

import (
	"context"
	"digital-book/internal/types"
)

type repo interface {
	CreateStudent(ctx context.Context, student *types.Student) error
	GetStudents(ctx context.Context) (res []types.Student, err error)
	GetStudent(ctx context.Context, id string) (res types.Student, err error)
	UpdateStudent(ctx context.Context, name, class, id string) error
	DeleteStudent(ctx context.Context, id string) error

	CreateSubject(ctx context.Context, subject *types.Subject) error
	GetSubjects(ctx context.Context) (res []types.Subject, err error)
	UpdateSubject(ctx context.Context, title, id string) error
	DeleteSubject(ctx context.Context, id string) error
	GetSubject(ctx context.Context, id string) (res types.Subject, err error)

	CreateMark(ctx context.Context, mark *types.Mark) error
	UpdateMark(ctx context.Context, studentID, subjectID, value, id string) error
	GetMarks(ctx context.Context) (res []types.MarkResponse, err error)
	GetMark(ctx context.Context, id string) (res types.Mark, err error)
	DeleteMark(ctx context.Context, id string) error
	AvgMarks(ctx context.Context) (res []types.MarkAverege, err error)
}

type srv struct {
	repo repo
}

func New(r repo) *srv {
	return &srv{repo: r}
}

package service

import (
	"context"
	"digital-book/internal/types"
)

type repo interface {
	CreateUser(ctx context.Context, user *types.User) error
	GetUsers(ctx context.Context) (res []types.User, err error)
	GetUser(ctx context.Context, id string) (res types.User, err error)
	UpdateUser(ctx context.Context, name, class, id, login string) error
	DeleteUser(ctx context.Context, id string) error

	CreateSubject(ctx context.Context, subject *types.Subject) error
	GetSubjects(ctx context.Context) (res []types.Subject, err error)
	UpdateSubject(ctx context.Context, title, id string) error
	DeleteSubject(ctx context.Context, id string) error
	GetSubject(ctx context.Context, id string) (res types.Subject, err error)

	CreateMark(ctx context.Context, mark *types.Mark) error
	UpdateMark(ctx context.Context, userID, subjectID, value, id string) error
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

package handlers

import (
	"context"
	"digital-book/internal/types"
	"net/http"
	"os"
	"text/template"
)

type repo interface {
	GetStudents(ctx context.Context) (res []types.Student, err error)
	GetStudent(ctx context.Context, id string) (res types.Student, err error)
	CreateStudent(ctx context.Context, name, class string) error
	UpdateStudent(ctx context.Context, name, class, id string) error
	DeleteStudent(ctx context.Context, id string) error

	CreateSubject(ctx context.Context, title string) error
	GetSubjects(ctx context.Context) (res []types.Subject, err error)
	UpdateSubject(ctx context.Context, title, id string) error
	DeleteSubject(ctx context.Context, id string) error

	CreateMark(ctx context.Context, studentID, subjectID, value string) error
	UpdateMark(ctx context.Context, studentID, subjectID, value, id string) error
	GetMarks(ctx context.Context) (res []types.MarkResponse, err error)
	GetMark(ctx context.Context, id string) (res types.Mark, err error)
	DeleteMark(ctx context.Context, id string) error
}

type server struct {
	repo repo
}

func New(r repo) *server {
	return &server{
		repo: r,
	}
}

func (*server) renderTemplate(fileName string, rw http.ResponseWriter, data any) {
	templBase, err := template.ParseFiles("html/base.html")
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	templ, err := os.ReadFile(fileName)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	_, err = templBase.Parse(string(templ))
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	templBase.Execute(rw, data)
}

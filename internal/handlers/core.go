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

	GetSubjects(ctx context.Context) (res []types.Subject, err error)
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

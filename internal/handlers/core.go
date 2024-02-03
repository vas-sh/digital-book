package handlers

import (
	"context"
	"digital-book/internal/types"
	"log"
	"net/http"
	"os"
	"text/template"
)

type service interface {
	CreateUser(ctx context.Context, user *types.User) error
	GetUsers(ctx context.Context) (res []types.User, err error)
	GetUser(ctx context.Context, id string) (res types.User, err error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, name, class, id, login string) error

	CreateSubject(ctx context.Context, subject *types.Subject) error
	GetSubject(ctx context.Context, id string) (res types.Subject, err error)
	GetSubjects(ctx context.Context) (res []types.Subject, err error)
	DeleteSubject(ctx context.Context, id string) error
	UpdateSubject(ctx context.Context, title, id string) error

	GetMarks(ctx context.Context) (res []types.MarkResponse, err error)
	GetMark(ctx context.Context, id string) (res types.Mark, err error)
	DeleteMark(ctx context.Context, id string) error
	CreateMark(ctx context.Context, mark *types.Mark) error
	UpdateMark(ctx context.Context, userID, subjectID, value, id string) error
	AvgMarks(ctx context.Context) (res []types.MarkAverege, err error)
}

type repo interface {
}

type server struct {
	repo repo
	srv  service
}

func New(r repo, s service) *server {
	return &server{
		repo: r,
		srv:  s,
	}
}

func (s *server) Run() {
	http.HandleFunc("/marks/avg", s.AvgMarks)

	http.HandleFunc("/marks", s.GetMarks)
	http.HandleFunc("/marks/create-new", s.CreateMarks)
	http.HandleFunc("/marks/delete", s.DeleteMark)

	http.HandleFunc("/subjects", s.GetSubjects)
	http.HandleFunc("/subjects/create-new", s.CreateSubject)
	http.HandleFunc("/subjects/delete", s.DeleteSubject)

	http.HandleFunc("/users/create-new", s.CreateUser)
	http.HandleFunc("/users", s.GetUsers)
	http.HandleFunc("/users/delete", s.DeleteUser)

	log.Println("start")
	http.ListenAndServe(":5005", nil)
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

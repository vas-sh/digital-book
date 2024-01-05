package handlers

import (
	"context"
	"digital-book/internal/types"
	"log"
	"net/http"
	"os"
	"text/template"
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

type server struct {
	repo repo
}

func New(r repo) *server {
	return &server{
		repo: r,
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

	http.HandleFunc("/students/create-new", s.CreateStudent)
	http.HandleFunc("/students", s.GetStudents)
	http.HandleFunc("/students/delete", s.DeleteStudent)

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

package handlers

import (
	"digital-book/internal/types"
	"net/http"
	"text/template"
)

func (s *server) GetSubjects(rw http.ResponseWriter, r *http.Request) {
	subjects, err := s.repo.GetSubjects(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	templ, err := template.ParseFiles("html/subjects.html")
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	data := struct {
		Subjects []types.Subject
	}{
		Subjects: subjects,
	}
	s.renderTemplate("html/subjects.html", rw, data)

	if err := templ.Execute(rw, data); err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

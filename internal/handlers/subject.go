package handlers

import (
	"digital-book/internal/types"
	"log"
	"net/http"
)

func (s *server) GetSubjects(rw http.ResponseWriter, r *http.Request) {
	subjects, err := s.repo.GetSubjects(r.Context())
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
}

func (s *server) CreateSubject(rw http.ResponseWriter, r *http.Request) {
	log.Println("createSubject", r.Method)
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost:
		title := r.FormValue("subject")
		id := r.FormValue("id")

		if id == "" || id == "0" {
			log.Println("new subject: title", title)
			if err := s.srv.CreateSubject(ctx, &types.Subject{Title: title}); err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		} else {
			log.Println("update subject: title", title, "id", id)
			if err := s.repo.UpdateSubject(ctx, title, id); err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		}

		http.Redirect(rw, r, "/subjects", http.StatusTemporaryRedirect)
		return

	case http.MethodGet:
		if id := r.URL.Query().Get("id"); id != "" {
			subject, err := s.srv.GetSubject(ctx, id)
			if err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}

			s.renderTemplate("html/update-subject.html", rw, struct {
				Subject types.Subject
			}{
				Subject: subject,
			})
		} else {
			s.renderTemplate("html/create-subject.html", rw, nil)
		}
	}
}

func (s *server) DeleteSubject(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		http.Error(rw, "not supported", http.StatusNotImplemented)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(rw, "id is required", http.StatusBadRequest)
		return
	}
	if err := s.repo.DeleteSubject(ctx, id); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/subjects", http.StatusTemporaryRedirect)
}

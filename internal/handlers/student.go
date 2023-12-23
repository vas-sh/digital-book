package handlers

import (
	"digital-book/internal/types"
	"log"
	"net/http"
)

func (s *server) GetStudents(rw http.ResponseWriter, r *http.Request) {
	students, err := s.repo.GetStudents(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	if len(students) == 0 {
		http.Error(rw, "no students", http.StatusNotFound)
		return
	}

	data := struct {
		Students []types.Student
	}{
		Students: students,
	}
	s.renderTemplate("html/students.html", rw, data)
}

func (s *server) CreateStudent(rw http.ResponseWriter, r *http.Request) {
	log.Println("createStudent", r.Method)
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost:
		class := r.FormValue("class")
		name := r.FormValue("name")
		id := r.FormValue("id")

		if id == "" {
			log.Println("new student: name", name, "class", class)
			if err := s.repo.CreateStudent(ctx, name, class); err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		} else {
			log.Println("update student: name", name, "class", class, "id", id)
			if err := s.repo.UpdateStudent(ctx, name, class, id); err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		}

		http.Redirect(rw, r, "/students", http.StatusTemporaryRedirect)
		return

	case http.MethodGet:
		if id := r.URL.Query().Get("id"); id != "" {
			student, err := s.repo.GetStudent(ctx, id)
			if err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
			s.renderTemplate("html/update-student.html", rw, struct {
				Student types.Student
			}{
				Student: student,
			})
		} else {
			s.renderTemplate("html/create-student.html", rw, nil)
		}
	}
}

func (s *server) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		http.Error(w, "not supported", http.StatusNotImplemented)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := s.repo.DeleteStudent(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/students", http.StatusTemporaryRedirect)
}

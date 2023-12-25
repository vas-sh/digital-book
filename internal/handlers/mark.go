package handlers

import (
	"digital-book/internal/types"
	"log"
	"net/http"
)

func (s *server) GetMarks(rw http.ResponseWriter, r *http.Request) {
	marks, err := s.repo.GetMarks(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	if len(marks) == 0 {
		http.Error(rw, "no marks", http.StatusNotFound)
		return
	}

	data := struct {
		Marks []types.MarkResponse
	}{
		Marks: marks,
	}

	s.renderTemplate("html/marks.html", rw, data)
}

func (s *server) DeleteMark(w http.ResponseWriter, r *http.Request) {
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
	if err := s.repo.DeleteMark(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/marks", http.StatusTemporaryRedirect)
}

func (s *server) CreateMarks(rw http.ResponseWriter, r *http.Request) {
	log.Println("createMark", r.Method)
	ctx := r.Context()
	switch r.Method {
	case http.MethodPost:
		studentID := r.FormValue("student_id")
		subjectID := r.FormValue("subject_id")
		value := r.FormValue("value")
		id := r.FormValue("id") ///
		if id == "" || id == "0" {
			log.Println("new mark: name", studentID, "lesson", subjectID, "point", value)
			if err := s.repo.CreateMark(ctx, studentID, subjectID, value); err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		} else {
			log.Println("update mark: student_id", studentID, "subject_id", subjectID, "value", value, "id", id)
			if err := s.repo.UpdateMark(ctx, studentID, subjectID, value, id); err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		}
		http.Redirect(rw, r, "/marks", http.StatusTemporaryRedirect)
		return

	case http.MethodGet:
		id := r.URL.Query().Get("id")

		subjects, err := s.repo.GetSubjects(ctx)
		if err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		students, err := s.repo.GetStudents(ctx)
		if err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		if id != "" {
			mark, err := s.repo.GetMark(ctx, id)
			if err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}

			s.renderTemplate("html/update-mark.html", rw, struct {
				Mark     types.Mark
				Students []types.Student
				Subjects []types.Subject
			}{
				Mark:     mark,
				Students: students,
				Subjects: subjects,
			})
		} else {
			s.renderTemplate("html/create-mark.html", rw, struct {
				Students []types.Student
				Subjects []types.Subject
			}{
				Students: students,
				Subjects: subjects,
			})
		}
	}
}

func (s *server) AvgMarks(rw http.ResponseWriter, r *http.Request) {

	avgMark, err := s.repo.AvgMarks(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	if len(avgMark) == 0 {
		http.Error(rw, "field is empty", http.StatusNotFound)
		return
	}

	data := struct {
		AvgMarks []types.MarkAverege
	}{
		AvgMarks: avgMark,
	}
	s.renderTemplate("html/AVG-marks.html", rw, data)

}

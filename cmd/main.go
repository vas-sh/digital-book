package main

import (
	"digital-book/internal/handlers"
	"digital-book/internal/repo"
	"digital-book/internal/types"
	"log"
	"net/http"
	"os"
	"text/template"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=user password=1111 dbname=test port=5432 sslmode=disable TimeZone=Europe/Kiev"

var db *gorm.DB

// TODO: move Mark, MarkReponse, MrAverege to the types/mark.go
type Mark struct {
	ID        int
	StudentID int
	SubjectID int
	Value     int
}

type MarkResponse struct {
	ID           int
	StudentName  string
	SubjectTitle string
	Value        int
}

type MrAverege struct {
	ID    int
	Name  string
	Title string
	Value float64
}

func renderTemplate(fileName string, rw http.ResponseWriter, data any) {
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

// TODO: move to handlers/mark.go
func AvgMarks(rw http.ResponseWriter, r *http.Request) {
	var avgMr []MrAverege
	if err := db.Raw(`
		SELECT student.id, student.name, subject.title, AVG(value) AS Value
	    FROM mark 
	        INNER JOIN student ON mark.student_id = student.id
	        INNER JOIN subject ON mark.subject_id = subject.id
	        GROUP BY student.id, student.name, subject.title 
			ORDER BY student.id ASC`).Scan(&avgMr).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	if len(avgMr) == 0 {
		http.Error(rw, "field is empty", http.StatusNotFound)
		return
	}

	templ, err := template.ParseFiles("html/AVG-marks.html")
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	data := struct {
		AvgMarks []MrAverege
	}{
		AvgMarks: avgMr,
	}
	renderTemplate("html/AVG-marks.html", rw, data)

	err = templ.Execute(rw, data)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

}

// TODO: move to handlers/mark.go
func createMarks(rw http.ResponseWriter, r *http.Request) {
	log.Println("createMark", r.Method)

	switch r.Method {
	case http.MethodPost:
		studentID := r.FormValue("student_id")
		subjectID := r.FormValue("subject_id")
		value := r.FormValue("value")
		id := r.FormValue("id")

		if id == "" || id == "0" {
			log.Println("new mark: name", studentID, "lesson", subjectID, "point", value)
			if err := db.Exec("INSERT INTO mark (ID, student_id, subject_id, value) VALUES(DEFAULT, ?, ?, ?)",
				studentID, subjectID, value).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		} else {
			log.Println("update mark: student_id", studentID, "subject_id", subjectID, "value", value, "id", id)
			if err := db.Exec("UPDATE mark SET student_id = ?, subject_id = ?, value = ? WHERE id = ?",
				studentID, subjectID, value, id).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}

		}
		http.Redirect(rw, r, "/marks", http.StatusTemporaryRedirect)
		return

	case http.MethodGet:
		var subjects []types.Subject
		if err := db.Raw("SELECT * FROM subject").Scan(&subjects).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		var students []types.Student
		if err := db.Raw("SELECT * FROM student").Scan(&students).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
		if id := r.URL.Query().Get("id"); id != "" {
			var mark Mark
			if err := db.Raw("SELECT * FROM mark WHERE id = ?", id).Scan(&mark).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
			renderTemplate("html/update-mark.html", rw, struct {
				Mark     Mark
				Students []types.Student
				Subjects []types.Subject
			}{
				Mark:     mark,
				Students: students, Subjects: subjects,
			})
		} else {

			renderTemplate("html/create-mark.html", rw, struct {
				Students []types.Student
				Subjects []types.Subject
			}{
				Students: students, Subjects: subjects,
			})
		}
	}
}

// TODO: move to handlers/mark.go
func getMarks(rw http.ResponseWriter, r *http.Request) {
	var marks []MarkResponse
	if err := db.Raw(`
		SELECT mark.id, student.name as student_name, subject.title as subject_title, mark.value 
		FROM mark
			 INNER JOIN student 
			 ON mark.student_id = student.id 
			 INNER JOIN subject 
			 ON mark.subject_id = subject.id
			 ORDER BY id ASC`).Scan(&marks).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	templ, err := template.ParseFiles("html/marks.html")
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	data := struct {
		Marks []MarkResponse
	}{
		Marks: marks,
	}
	renderTemplate("html/marks.html", rw, data)

	err = templ.Execute(rw, data)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
}

// TODO: move handlers/subject.go

// TODO: move to handlers/mark.go
func deleteMark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not supported", http.StatusNotImplemented)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := db.Exec("DELETE FROM mark WHERE id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/marks", http.StatusTemporaryRedirect)
}

// TODO: Move to handlers/student.go
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not supported", http.StatusNotImplemented)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := db.Exec("DELETE FROM student WHERE id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/students", http.StatusTemporaryRedirect)
}

// TODO: move to handlers/subject.go
func deleteSubject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not supported", http.StatusNotImplemented)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := db.Exec("DELETE FROM subject WHERE id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/subjects", http.StatusTemporaryRedirect)
}

func main() {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rep := repo.New(db)
	server := handlers.New(rep)

	http.HandleFunc("/marks/avg", AvgMarks)

	http.HandleFunc("/marks", getMarks)
	http.HandleFunc("/marks/create-new", createMarks)
	http.HandleFunc("/marks/delete", deleteMark)

	http.HandleFunc("/subjects", server.GetSubjects)
	http.HandleFunc("/subjects/create-new", server.CreateSubject)
	http.HandleFunc("/subjects/delete", deleteSubject)

	http.HandleFunc("/students/create-new", server.CreateStudent)
	http.HandleFunc("/students", server.GetStudents)
	http.HandleFunc("/students/delete", deleteStudent)

	log.Println("start")
	http.ListenAndServe(":5005", nil)
}

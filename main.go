package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=user password=1111 dbname=test port=5432 sslmode=disable TimeZone=Europe/Kiev"

var db *gorm.DB

type Subject struct {
	ID    int
	Title string
}

type Student struct {
	ID    int
	Name  string
	Class string
}

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
		var subjects []Subject
		if err := db.Raw("SELECT * FROM subject").Scan(&subjects).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		var students []Student
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
				Students []Student
				Subjects []Subject
			}{
				Mark:     mark,
				Students: students, Subjects: subjects,
			})
		} else {

			renderTemplate("html/create-mark.html", rw, struct {
				Students []Student
				Subjects []Subject
			}{
				Students: students, Subjects: subjects,
			})
		}
	}
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

func getStudents(rw http.ResponseWriter, r *http.Request) {
	var students []Student
	if err := db.Raw("SELECT * FROM student ORDER BY id ASC").Scan(&students).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	if len(students) == 0 {
		http.Error(rw, "no students", http.StatusNotFound)
		return
	}

	data := struct {
		Students []Student
	}{
		Students: students,
	}
	renderTemplate("html/students.html", rw, data)

}

func getSubjects(rw http.ResponseWriter, r *http.Request) {
	var subjects []Subject
	if err := db.Raw("SELECT * FROM subject ORDER BY id").Scan(&subjects).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	templ, err := template.ParseFiles("html/subjects.html")
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	data := struct {
		Subjects []Subject
	}{
		Subjects: subjects,
	}
	renderTemplate("html/subjects.html", rw, data)

	if err := templ.Execute(rw, data); err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

func createStudent(rw http.ResponseWriter, r *http.Request) {
	log.Println("createStudent", r.Method)

	switch r.Method {
	case http.MethodPost:
		class := r.FormValue("class")
		name := r.FormValue("name")
		id := r.FormValue("id")

		if id == "" {
			log.Println("new student: name", name, "class", class)
			if err := db.Exec("INSERT INTO student (name, class) VALUES (?, ?)", name, class).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		} else {
			log.Println("update student: name", name, "class", class, "id", id)
			if err := db.Exec("UPDATE student SET name = ?, class = ? WHERE id = ?", name, class, id).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		}

		http.Redirect(rw, r, "/students", http.StatusTemporaryRedirect)
		return

	case http.MethodGet:
		if id := r.URL.Query().Get("id"); id != "" {
			var student Student
			if err := db.Raw("SELECT * FROM student WHERE id = ?", id).Scan(&student).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
			renderTemplate("html/update-student.html", rw, struct {
				Student Student
			}{
				Student: student,
			})
		} else {
			renderTemplate("html/create-student.html", rw, nil)
		}
	}
}

func createSubject(rw http.ResponseWriter, r *http.Request) {
	log.Println("createSubject", r.Method)

	switch r.Method {
	case http.MethodPost:

		title := r.FormValue("subject")
		id := r.FormValue("id")

		if id == "" || id == "0" {
			log.Println("new subject: title", title)
			if err := db.Exec("INSERT INTO subject (ID, title) VALUES(DEFAULT, ?)",
				title).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		} else {
			log.Println("updade subject: title", title, id)
			if err := db.Exec("UPDATE subject SET title = ? WHERE id = ?",
				title, id).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		}

		http.Redirect(rw, r, "/subjects", http.StatusTemporaryRedirect)
		return

	case http.MethodGet:

		if id := r.URL.Query().Get("id"); id != "" {
			var subject Subject
			if err := db.Raw("SELECT * FROM subject WHERE id = ?", id).Scan(&subject).Error; err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
			renderTemplate("html/update-subject.html", rw, struct {
				Subject Subject
			}{
				Subject: subject,
			})
		} else {
			renderTemplate("html/create-subject.html", rw, nil)
		}
	}
}

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
	http.HandleFunc("/marks/avg", AvgMarks)

	http.HandleFunc("/marks", getMarks)
	http.HandleFunc("/marks/create-new", createMarks)
	http.HandleFunc("/marks/delete", deleteMark)

	http.HandleFunc("/subjects", getSubjects)
	http.HandleFunc("/subjects/create-new", createSubject)
	http.HandleFunc("/subjects/delete", deleteSubject)

	http.HandleFunc("/students/create-new", createStudent)
	http.HandleFunc("/students", getStudents)
	http.HandleFunc("/students/delete", deleteStudent)

	log.Println("start")
	http.ListenAndServe(":5005", nil)
}

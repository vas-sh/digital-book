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
	if r.Method == http.MethodPost {
		name := r.FormValue("student_id")
		lesson := r.FormValue("subject_id")
		point := r.FormValue("value")

		log.Println("new mark: student_id", name, "subject_id", lesson, "value", point)

		if err := db.Exec("INSERT INTO mark(student_id, subject_id, value) VALUES(?, ?, ?)",
			name, lesson, point).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
		getMarks(rw, r)
		return
	}

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

	renderTemplate("html/create-mark.html", rw, struct {
		Students []Student
		Subjects []Subject
	}{
		Students: students, Subjects: subjects,
	})
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
			 ON mark.subject_id = subject.id`).Scan(&marks).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	if len(marks) == 0 {
		http.Error(rw, "field is empty", http.StatusNotFound)
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
	if err := db.Raw("SELECT * FROM student").Scan(&students).Error; err != nil {
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
	if err := db.Raw("SELECT * FROM subject").Scan(&subjects).Error; err != nil {
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
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		class := r.FormValue("class")

		log.Println("new student: name", name, "class", class)

		if err := db.Exec("INSERT INTO student(name, class) VALUES(?, ?)",
			name, class).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
		getStudents(rw, r)
		return
	}
	renderTemplate("html/create-student.html", rw, nil)
}

func createSubject(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("subject")

		log.Println("new subject: title", title)

		if err := db.Exec("INSERT INTO subject (ID, title) VALUES(DEFAULT, ?)",
			title).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
		getSubjects(rw, r)
		return
	} else if r.Method == http.MethodPatch {
		title := r.FormValue("subject")
		id := r.FormValue("id")
		log.Println("updade subject: title", title, id)

		if err := db.Exec("UPDATE subject SET title = ? WHERE id = ?",
			title, id).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
		getSubjects(rw, r)
		return
	}
	id := r.URL.Query().Get("id")
	var subject Subject
	var method = http.MethodPost
	if id != "" {
		if err := db.Raw("SELECT * FROM subject WHERE id = ?", id).Scan(&subject).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
		method = http.MethodPatch
	}
	renderTemplate("html/create-subject.html", rw, struct {
		Subject Subject
		Method  string
	}{
		Subject: subject, Method: method,
	})
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

	http.HandleFunc("/subjects", getSubjects)
	http.HandleFunc("/subjects/create-new", createSubject)

	http.HandleFunc("/students/create-new", createStudent)
	http.HandleFunc("/students", getStudents)

	log.Println("start")
	http.ListenAndServe(":5005", nil)
}

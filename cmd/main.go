package main

import (
	"digital-book/internal/handlers"
	"digital-book/internal/repo"

	//"digital-book/internal/types"
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

type MarkAverege struct {
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
	var avgMr []MarkAverege
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
		AvgMarks []MarkAverege
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

// TODO: move to handlers/mark.go

func main() {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rep := repo.New(db)
	server := handlers.New(rep)

	http.HandleFunc("/marks/avg", AvgMarks)

	http.HandleFunc("/marks", server.GetMarks)
	http.HandleFunc("/marks/create-new", server.CreateMarks)
	http.HandleFunc("/marks/delete", server.DeleteMark)

	http.HandleFunc("/subjects", server.GetSubjects)
	http.HandleFunc("/subjects/create-new", server.CreateSubject)
	http.HandleFunc("/subjects/delete", server.DeleteSubject)

	http.HandleFunc("/students/create-new", server.CreateStudent)
	http.HandleFunc("/students", server.GetStudents)
	http.HandleFunc("/students/delete", server.DeleteStudent)

	log.Println("start")
	http.ListenAndServe(":5005", nil)
}

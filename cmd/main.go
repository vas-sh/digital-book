package main

import (
	"digital-book/internal/handlers"
	"digital-book/internal/repo"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=user password=1111 dbname=test port=5432 sslmode=disable TimeZone=Europe/Kiev"

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rep := repo.New(db)
	server := handlers.New(rep)

	http.HandleFunc("/marks/avg", server.AvgMarks)

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

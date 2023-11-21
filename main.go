package main

import (
	"fmt"
	"log"
	"net/http"

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

func createMarks(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("Name")
		lesson := r.FormValue("subject_id")
		point := r.FormValue("value")

		log.Println("new student: name", name, "class", lesson)

		if err := db.Exec("INSERT INTO mark(student_id, subject_id, value) VALUES(?, ?, ?)",
			name, lesson, point).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
	}
	var subjects []Subject
	if err := db.Raw("SELECT * FROM subject").Scan(&subjects).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	var options string
	for _, subject := range subjects {
		options += fmt.Sprintf(`<option value="%d">%s</option>`, subject.ID, subject.Title)
	}

	var students []Student
	if err := db.Raw("SELECT * FROM student").Scan(&students).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	var stdOptions string
	for _, student := range students {
		stdOptions += fmt.Sprintf(`<option value="%d">%s</option>`, student.ID, student.Name)
	}

	rw.Write([]byte(`
	<html>
<body>
	<form method="POST" action="/marks/create-new">
	<label>Name: </label>
	<select name="Name">
 ` + stdOptions + `
</select>
   <br>
      
	  <label>Lesson: </label>
       <select name="subject_id">
    ` + options + `
  </select>
      <br>
      
	  <label>Value: </label>
      <input required name="value" type="number" />
      <br>
      
	  <button type="submit">Create</button>
	</form>
 
</body>
</html>

	`))
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
	text := ""
	if len(marks) == 0 {
		text = "no students"
	}
	for _, mark := range marks {
		text += fmt.Sprintf("<p>%d. NameStudent: %s, subject: %s, value: %d</p>", mark.ID, mark.StudentName, mark.SubjectTitle, mark.Value)
	}
	rw.Write([]byte(text))
}

func getStudents(rw http.ResponseWriter, r *http.Request) {
	var students []Student
	if err := db.Raw("SELECT * FROM student").Scan(&students).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	text := ""
	if len(students) == 0 {
		text = "no students"
	}
	for _, student := range students {
		text += fmt.Sprintf("<p>%d. Name: %s, Class: %s</p>", student.ID, student.Name, student.Class)
	}
	rw.Write([]byte(text))
}

func getSubjects(rw http.ResponseWriter, r *http.Request) {
	var subjects []Subject
	if err := db.Raw("SELECT * FROM subject").Scan(&subjects).Error; err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	text := ""
	if len(subjects) == 0 {
		text = "no students"
	}
	for _, subject := range subjects {
		text += fmt.Sprintf("<p>%d Subject: %s</p>", subject.ID, subject.Title)
	}
	rw.Write([]byte(text))
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
	}
	rw.Write([]byte(`
	 <body>
	  <form method="POST">
	   <label>Name: </label><input required name="name" type="text" /></br>
	   <label>Class: </label><input required name="class" type="number" /></br>
	   <button action="submit">Create</button>
	  </form>
	 </body>
	`))
}

func createSubject(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")

		log.Println("new subject: title", title)

		if err := db.Exec("INSERT INTO subject (title) VALUES(?)",
			title).Error; err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
	}
	rw.Write([]byte(`
	 <body>
	  <form method="POST">
	   <label>Title: </label><input required name="title" type="text" /></br>
	   <button action="submit">Create</button>
	  </form>
	 </body>
	`))
}

func main() {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/marks", getMarks)
	http.HandleFunc("/marks/create-new", createMarks)

	http.HandleFunc("/subjects", getSubjects)
	http.HandleFunc("/subjects/create-new", createSubject)

	http.HandleFunc("/students/create-new", createStudent)
	http.HandleFunc("/students", getStudents)

	log.Println("start")
	http.ListenAndServe(":5005", nil)
}

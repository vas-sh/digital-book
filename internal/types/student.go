package types

// TODO: add login and password, token fields (db type is string)
// Make login and token fields unique
// Rename the student to user and all the related fields from student_id to user_id
type Student struct {
	ID    int `gorm:"primary_key"`
	Name  string
	Class string
	Login string `gorm:"unique"`
}

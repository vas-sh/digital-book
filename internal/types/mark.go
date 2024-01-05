package types

type Mark struct {
	ID        int `gorm:"primary_key"`
	StudentID int
	SubjectID int
	Value     int
}

type MarkResponse struct {
	ID           int `gorm:"primary_key"`
	StudentName  string
	SubjectTitle string
	Value        int
}

type MarkAverege struct {
	ID    int `gorm:"primary_key"`
	Name  string
	Title string
	Value float64
}

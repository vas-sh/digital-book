package types

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

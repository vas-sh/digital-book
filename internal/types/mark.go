package types

type Mark struct {
	ID        int
	UserID    int
	SubjectID int
	Value     int
}

type MarkResponse struct {
	ID           int
	UserName     string
	SubjectTitle string
	Value        int
}

type MarkAverege struct {
	ID    int
	Name  string
	Title string
	Value float64
}

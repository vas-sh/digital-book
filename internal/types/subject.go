package types

type Subject struct {
	ID    int `gorm:"primary_key"`
	Title string
}

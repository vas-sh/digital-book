package types

type Student struct {
	ID    int `gorm:"primary_key"`
	Name  string
	Class string
}

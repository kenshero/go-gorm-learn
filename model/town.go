package model

type Town struct {
	ID   int `gorm:"primary_key"`
	Name string
}

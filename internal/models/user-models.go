package models

import "gorm.io/gorm"

type Role string

const (
	Admin  = "Admin"
	Player = "player"
)

type User struct {
	*gorm.Model
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
}

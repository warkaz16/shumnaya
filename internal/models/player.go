package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Name         string
	Email        string
	PasswordHash string
	Rating       int
}

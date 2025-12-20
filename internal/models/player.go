package models

import "gorm.io/gorm"

// Player represents a player in the system.
type Player struct {
	gorm.Model
	Name         string
	Email        string
	PasswordHash string
	Rating       int
}

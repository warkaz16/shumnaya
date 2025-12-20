package models

import "gorm.io/gorm"

// Player represents a participant in the system.
type Player struct {
	gorm.Model
	Name         string
	Email        string
	PasswordHash string
	Rating       int

	// Relations
	Wins      []Match    `gorm:"foreignKey:WinnerID"`
	Losses    []Match    `gorm:"foreignKey:LoserID"`
	Standings []Standing `gorm:"foreignKey:PlayerID"`
}

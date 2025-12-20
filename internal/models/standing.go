package models

import "gorm.io/gorm"

// Standing represents a player's standing within a season.
type Standing struct {
	gorm.Model
	PlayerID uint
	SeasonID uint
	Wins     int
	Losses   int
	Points   int
	Rank     int
}

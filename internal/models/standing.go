package models

import "gorm.io/gorm"

// Standing represents a player's record within a season.
type Standing struct {
	gorm.Model

	PlayerID uint
	Player   Player

	SeasonID uint
	Season   Season

	Wins   int
	Losses int
	Points int
	Rank   int
}

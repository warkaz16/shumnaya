package models

import (
	"time"

	"gorm.io/gorm"
)

// Match represents a played game between two players within a season.
type Match struct {
	gorm.Model

	WinnerID uint
	Winner   Player

	LoserID uint
	Loser   Player

	SeasonID uint
	Season   Season

	Score              string
	WinnerRatingChange int
	LoserRatingChange  int
	PlayedAt           time.Time
}

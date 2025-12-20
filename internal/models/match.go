package models

import (
	"time"

	"gorm.io/gorm"
)

// Match represents a played match between two players.
type Match struct {
	gorm.Model
	WinnerID           uint
	LoserID            uint
	SeasonID           uint
	Score              string
	WinnerRatingChange int
	LoserRatingChange  int
	PlayedAt           time.Time
}

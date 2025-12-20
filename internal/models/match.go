package models

import (
	"time"

	"gorm.io/gorm"
)

// Match represents a played match between two players.
type Match struct {
	gorm.Model
	WinnerID           uint
	Winner             *Player  `gorm:"foreignKey:WinnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LoserID            uint
	Loser              *Player  `gorm:"foreignKey:LoserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SeasonID           uint
	Season             *Season  `gorm:"foreignKey:SeasonID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Score              string
	WinnerRatingChange int
	LoserRatingChange  int
	PlayedAt           time.Time
}

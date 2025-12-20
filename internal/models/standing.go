package models

import "gorm.io/gorm"

// Standing represents a player's standing within a season.
type Standing struct {
	gorm.Model
	PlayerID uint
	Player   *Player `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SeasonID uint
	Season   *Season `gorm:"foreignKey:SeasonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Wins     int
	Losses   int
	Points   int
	Rank     int
}

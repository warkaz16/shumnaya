package models

import "gorm.io/gorm"

type Standing struct {
	gorm.Model

	PlayerID uint
	Player   Player

	SeasonID uint
	Season   Season

	Wins     int
	Losses   int
	Points   int
	Rank     int
}

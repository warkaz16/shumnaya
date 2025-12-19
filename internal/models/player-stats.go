package models

import "gorm.io/gorm"

type PlayerStats struct {
	gorm.Model

	UserID      uint
	SeasonID    uint

	GamesPlayed int
	Wins        int
	Loses       int
	Elo         int
}
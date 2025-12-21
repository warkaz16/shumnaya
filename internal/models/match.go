package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model

	WinnerID            uint
	Winner              Player

	LoserID             uint
	Loser               Player

	SeasonID            uint
	Season              Season

	Score               string
	WinnerRatingChange  int
	LoserRatingChange   int
	PlayedAt            time.Time
}


type MatchFilter struct {
	SeasonID *uint
	PlayerID *uint
	FromDate *time.Time
	ToDate   *time.Time
}
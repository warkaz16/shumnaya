package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model

	WinnerID uint   `json:"winner_id" gorm:"column:winner_id"`
	Winner   Player `json:"winner,omitempty" gorm:"foreignKey:WinnerID;references:ID"`

	LoserID uint   `json:"loser_id" gorm:"column:loser_id"`
	Loser   Player `json:"loser,omitempty" gorm:"foreignKey:LoserID;references:ID"`

	SeasonID uint   `json:"season_id" gorm:"column:season_id"`
	Season   Season `json:"season,omitempty" gorm:"foreignKey:SeasonID;references:ID"`

	Score              string    `json:"score" gorm:"column:score"`
	WinnerRatingChange int       `json:"winner_rating_change" gorm:"column:winner_rating_change"`
	LoserRatingChange  int       `json:"loser_rating_change" gorm:"column:loser_rating_change"`
	PlayedAt           time.Time `json:"played_at" gorm:"column:played_at"`
}

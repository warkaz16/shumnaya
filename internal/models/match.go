package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model

	WinnerID uint   `json:"winner_id" gorm:"column:winner_id" binding:"required,min=1"`
	Winner   Player `json:"winner,omitempty" gorm:"foreignKey:WinnerID;references:ID"`

	LoserID uint   `json:"loser_id" gorm:"column:loser_id" binding:"required,min=1"`
	Loser   Player `json:"loser,omitempty" gorm:"foreignKey:LoserID;references:ID"`

	SeasonID uint   `json:"season_id" gorm:"column:season_id" binding:"required,min=1"`
	Season   Season `json:"season,omitempty" gorm:"foreignKey:SeasonID;references:ID"`

	Score              string    `json:"score" gorm:"column:score" binding:"required"`
	WinnerRatingChange int       `json:"winner_rating_change,omitempty" gorm:"column:winner_rating_change"`
	LoserRatingChange  int       `json:"loser_rating_change,omitempty" gorm:"column:loser_rating_change"`
	PlayedAt           time.Time `json:"played_at" gorm:"column:played_at"`
}

type MatchFilter struct {
	SeasonID *uint      `json:"season_id"`
	PlayerID *uint      `json:"player_id"`
	FromDate *time.Time `json:"from_date"`
	ToDate   *time.Time `json:"to_date"`
}

type CreateMatchRequest struct {
	WinnerID uint   `json:"winner_id" binding:"required"`
	LoserID  uint   `json:"loser_id" binding:"required"`
	SeasonID uint   `json:"season_id" binding:"required"`
	Score    string `json:"score" binding:"required"`
}

type HeadToHeadRecord struct {
	PlayerAID         uint      `json:"player_a_id"`
	PlayerBID         uint      `json:"player_b_id"`
	PlayerAWins       int      `json:"player_a_wins"`
	PlayerBWins       int      `json:"player_b_wins"`
	TotalMatches      int      `json:"total_matches"`
	LastMatchesPlayed []Match   `json:"last_matches_played"`
}

package models

import "gorm.io/gorm"

type Standing struct {
	gorm.Model

	PlayerID uint   `json:"player_id" binding:"required,min=1"`
	Player   Player `json:"player"`

	SeasonID uint   `json:"season_id" binding:"required,min=1"`
	Season   Season `json:"season"`

	Wins   int `json:"wins" binding:"min=0"`
	Losses int `json:"losses" binding:"min=0"`
	Points int `json:"points" binding:"min=0"`
	Rank   int `json:"rank" binding:"min=0"`
}

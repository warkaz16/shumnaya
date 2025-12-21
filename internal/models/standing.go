package models

import "gorm.io/gorm"

type Standing struct {
	gorm.Model

	PlayerID uint   `json:"player_id" gorm:"column:player_id"`
	Player   Player `json:"player,omitempty" gorm:"foreignKey:PlayerID;references:ID"`

	SeasonID uint   `json:"season_id" gorm:"column:season_id"`
	Season   Season `json:"season,omitempty" gorm:"foreignKey:SeasonID;references:ID"`

	Wins   int `json:"wins" gorm:"column:wins"`
	Losses int `json:"losses" gorm:"column:losses"`
	Points int `json:"points" gorm:"column:points"`
	Rank   int `json:"rank" gorm:"column:rank"`
}

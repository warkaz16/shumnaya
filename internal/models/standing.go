package models

import "gorm.io/gorm"

type Standing struct {
	gorm.Model `json:"-"`

	PlayerID uint   `json:"player_id" gorm:"column:player_id" binding:"required,min=1"`
	Player   Player `json:"player,omitempty" gorm:"foreignKey:PlayerID;references:ID"`

	SeasonID uint   `json:"season_id" gorm:"column:season_id" binding:"required,min=1"`
	Season   Season `json:"season,omitempty" gorm:"foreignKey:SeasonID;references:ID"`

	Wins   int `json:"wins" gorm:"column:wins" binding:"min=0"`
	Losses int `json:"losses" gorm:"column:losses" binding:"min=0"`
	Points int `json:"points" gorm:"column:points" binding:"min=0"`
	Rank   int `json:"rank" gorm:"column:rank" binding:"min=0"`
}

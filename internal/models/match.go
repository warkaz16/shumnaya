package models

import "gorm.io/gorm"

type Match struct {
	*gorm.Model

	SeasonID uint
	Season   Season

	Player1ID uint
	Player1   User

	Player2ID uint
	Player2   User

	Sets []MatchSet
}

type MatchSet struct {
	*gorm.Model

	MatchID      uint
	SetNumber    int
	Player1Score uint
	Player2Score uint
}

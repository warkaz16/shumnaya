package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Name         string  `json:"name" gorm:"column:name;type:varchar(255)" binding:"required"`
	Email        string  `json:"email" gorm:"column:email;type:varchar(255);uniqueIndex" binding:"required,email"`
	PasswordHash string  `json:"password_hash,omitempty" gorm:"column:password_hash"`
	Rating       int     `json:"rating" gorm:"column:rating" binding:"min=0"`

	Matches []Match `json:"matches,omitempty" gorm:"-"` // история матчей по игроку (поле для удобства, запросы через репозиторий)
}


type PlayerProfile struct {
	Player        Player  `json:"player"`
	Rating        int     `json:"rating"`
	TotalMatches  int     `json:"total_matches"`
	Wins          int     `json:"wins"`
	Losses        int     `json:"losses"`
	RecentMatches []Match `json:"recent_matches,omitempty"`
}

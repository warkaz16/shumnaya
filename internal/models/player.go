package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Name         string  `json:"name" gorm:"column:name;type:varchar(255)"`
	Email        string  `json:"email" gorm:"column:email;type:varchar(255);uniqueIndex"`
	PasswordHash string  `json:"password_hash" gorm:"column:password_hash"`
	Rating       int     `json:"rating" gorm:"column:rating"`

	Matches []Match `json:"matches,omitempty" gorm:"-"` // история матчей по игроку (поле для удобства, запросы через репозиторий)
}

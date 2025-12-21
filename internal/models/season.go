package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	Name      string    `json:"name" gorm:"column:name;type:varchar(255)" binding:"required"`
	StartDate time.Time `json:"start_date" gorm:"column:start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" gorm:"column:end_date" binding:"required"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active"`

	Matches []Match `json:"matches,omitempty" gorm:"foreignKey:SeasonID"` // получение матчей по сезонам
}

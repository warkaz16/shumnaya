package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	Name      string     `json:"name" gorm:"not null"`
	IsActive  bool       `json:"is_active" gorm:"default:false"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	Matches []Match
	Stats   []PlayerStats
}

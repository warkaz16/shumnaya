package models

import (
	"time"

	"gorm.io/gorm"
)

// Season represents a competition period.
type Season struct {
	gorm.Model
	Name      string
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool

	// Relations
	Matches   []Match    `gorm:"foreignKey:SeasonID"`
	Standings []Standing `gorm:"foreignKey:SeasonID"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// Season represents a season for matches and standings.
type Season struct {
	gorm.Model
	Name      string
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	Name      string
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
}

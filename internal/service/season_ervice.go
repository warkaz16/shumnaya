package service

import "shumnaya/internal/models"

type SeasonService interface {
	Create(season *models.Season) (*models.Season, error)
	Get() ([]models.Season, error)
	GetByID(id uint) (*models.Season, error)
	Update(id uint) error
	Delete(id uint) error
}

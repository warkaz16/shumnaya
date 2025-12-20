package repository

import "shumnaya/internal/models"

type SeasonRepository interface {
	Create(season *models.Season) error

	GetByID(id uint) (*models.Season, error)
	GetActive() (*models.Season, error)
	GetAll() ([]models.Season, error)

	Update(season *models.Season) error
	CloseSeason(id uint) error
}

package repository

import "shumnaya/internal/models"

type StandingRepository interface {
	Create(standing *models.Standing) error
	Update(standing *models.Standing) error

	CreateOrUpdate(standing *models.Standing) error

	GetByPlayerAndSeason(playerID, seasonID uint) (*models.Standing, error)
	GetBySeason(seasonID uint) ([]models.Standing, error)
}

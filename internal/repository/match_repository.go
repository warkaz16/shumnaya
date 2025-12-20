package repository

import "shumnaya/internal/models"

type MatchRepository interface {
	Create(match *models.Match) error

	GetByID(id uint) (*models.Match, error)

	GetBySeasonID(seasonID uint) ([]models.Match, error)
	GetByPlayerID(playerID uint) ([]models.Match, error)
}

package service

import "shumnaya/internal/models"

type MatchService interface {
	Create(match *models.Match) (*models.Match, error)
	Get() ([]models.Match, error)
	GetByID(id uint) (*models.Match, error)
	Update(id uint) error
	Delete(id uint) error
}

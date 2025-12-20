package repository

import "shumnaya/internal/models"

type MatchRepo interface {
	Get() ([]models.Match, error)
	GetByID(ID uint) (*models.Match, error)
	Delete(ID uint) error
	Update(ID uint) error
	Create(*models.Match) error
}

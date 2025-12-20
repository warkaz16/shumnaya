package repository

import "shumnaya/internal/models"

type SeasonRepo interface {
	Create(*models.Season)error
	Get()([]models.Season, error)
	GetByID(ID uint)(*models.Season, error)
	Delete(ID uint) error
	Update(ID uint)error
}

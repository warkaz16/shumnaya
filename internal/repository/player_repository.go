package repository

import "shumnaya/internal/models"

type PlayerRepository interface {
	Create(player *models.Player) error

	GetByID(id uint) (*models.Player, error)
	GetByEmail(email string) (*models.Player, error)

	Update(player *models.Player) error
	Delete(id uint) error
}

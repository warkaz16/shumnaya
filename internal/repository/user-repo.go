package repository

import "shumnaya/internal/models"

type UserRepo interface {
	Get() ([]models.User, error)
	GetByID(ID uint) (*models.User, error)
	Delete(ID uint) error
	Create(*models.User) error
	Update(ID uint) error
}

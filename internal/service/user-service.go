package service

import "shumnaya/internal/models"

type UserService interface {
	Create(user *models.User) (*models.User, error)
	Get() ([]models.User, error)
	GetByID(id uint) (*models.User, error)
	Update(id uint) error
	Delete(id uint) error
}

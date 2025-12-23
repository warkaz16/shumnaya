package repository

import (
	"log/slog"
	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type PlayerRepository interface {
    WithDB(tx *gorm.DB) PlayerRepository
	Create(player *models.Player) error

	GetByID(id uint) (*models.Player, error)
	GetByEmail(email string) (*models.Player, error)

	Update(player *models.Player) error
	Delete(id uint) error
}

type playerRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewPlayerRepository(db *gorm.DB, logger *slog.Logger) PlayerRepository {
	return &playerRepository{db: db, logger: logger}
}

func (r *playerRepository) WithDB(tx *gorm.DB) PlayerRepository {
	return &playerRepository{db: tx, logger: r.logger}
}

func (r *playerRepository) GetByID(id uint) (*models.Player, error) {
	var player models.Player

	err := r.db.First(&player, id).Error
	if err != nil {
		r.logger.Error("ошибка получения игрока по ID", "id", id, "error", err)
		return nil, err
	}

	return &player, nil
}

func (r *playerRepository) GetByEmail(email string) (*models.Player, error) {
	var player models.Player

	err := r.db.Where("email = ?", email).First(&player).Error
	if err != nil {
		r.logger.Error("ошибка получения игрока по email", "email", email, "error", err)
		return nil, err
	}

	return &player, nil
}

func (r *playerRepository) Create(player *models.Player) error {
	err := r.db.Create(player).Error
	if err != nil {
		r.logger.Error("ошибка создания игрока", "error", err)
		return err
	}
	return nil
}

func (r *playerRepository) Update(player *models.Player) error {
	err := r.db.Save(player).Error
	if err != nil {
		r.logger.Error("ошибка обновления игрока", "id", player.ID, "error", err)
		return err
	}
	return nil
}

func (r *playerRepository) Delete(id uint) error {
	err := r.db.Delete(&models.Player{}, id).Error
	if err != nil {
		r.logger.Error("ошибка удаления игрока", "id", id, "error", err)
		return err
	}

	return nil
}

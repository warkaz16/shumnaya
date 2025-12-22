package repository

import (
	"log/slog"

	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type SeasonRepository interface {
	Create(season *models.Season) error

	GetByID(id uint) (*models.Season, error)
	GetActive() (*models.Season, error)
	GetAll() ([]models.Season, error)

	Update(season *models.Season) error
	CloseSeason(id uint) error
}

type seasonRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewSeasonRepository(db *gorm.DB, logger *slog.Logger) SeasonRepository {
	return &seasonRepository{
		db:     db,
		logger: logger,
	}
}

func (r *seasonRepository) Create(season *models.Season) error {
	if r.logger != nil {
		r.logger.Info("создание сезона", "name", season.Name)
	}

	if err := r.db.Create(season).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("ошибка при создании сезона", "error", err)
		}
		return err
	}

	if r.logger != nil {
		r.logger.Info("сезон успешно создан", "season_id", season.ID)
	}

	return nil
}

func (r *seasonRepository) GetAll() ([]models.Season, error) {
	var seasons []models.Season

	if err := r.db.Find(&seasons).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("ошибка при получении списка сезонов", "error", err)
		}
		return nil, err
	}
	if r.logger != nil {
		r.logger.Info("список сезонов получен", "count", len(seasons))
	}
	return seasons, nil
}

func (r *seasonRepository) GetByID(id uint) (*models.Season, error) {
	var season models.Season

	if err := r.db.First(&season, id).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("ошибка при получении сезона по id", "season_id", id, "error", err)
		}
		return nil, err
	}

	return &season, nil
}

func (r *seasonRepository) GetActive() (*models.Season, error) {
	var season models.Season

	if r.logger != nil {
		r.logger.Info("получение активного сезона")
	}

	if err := r.db.Where("is_active = ?", true).First(&season).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("активный сезон не найден", "error", err)
		}
		return nil, err
	}

	return &season, nil
}

func (r *seasonRepository) Update(season *models.Season) error {
	if err := r.db.Save(season).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("ошибка при обновлении сезона", "season_id", season.ID, "error", err)
		}
		return err
	}

	if r.logger != nil {
		r.logger.Info("сезон обновлён", "season_id", season.ID)
	}

	return nil
}

func (r *seasonRepository) CloseSeason(id uint) error {
	var season models.Season

	if err := r.db.First(&season, id).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("ошибка при закрытии сезона: сезон не найден", "season_id", id, "error", err)
		}
		return err
	}

	season.IsActive = false

	if err := r.db.Save(&season).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("ошибка при сохранении закрытого сезона", "season_id", id, "error", err)
		}
		return err
	}

	if r.logger != nil {
		r.logger.Info("сезон успешно закрыт", "season_id", id)
	}

	return nil
}

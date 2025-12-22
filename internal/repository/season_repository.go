package repository

import (
	"log/slog"

	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type SeasonRepository interface {
	Create(season *models.Season) error
	GetByID(id uint) (*models.Season, error)
	GetAll() ([]models.Season, error)
	GetActive() (*models.Season, error)
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
	if err := r.db.Create(season).Error; err != nil {
		if r.logger != nil {
			r.logger.Error(
				"repository: ошибка создания сезона",
				"season_name", season.Name,
				"error", err,
			)
		}
		return err
	}
	return nil
}

func (r *seasonRepository) GetAll() ([]models.Season, error) {
	var seasons []models.Season

	if err := r.db.Find(&seasons).Error; err != nil {
		if r.logger != nil {
			r.logger.Error(
				"repository: ошибка получения списка сезонов",
				"error", err,
			)
		}
		return nil, err
	}
	return seasons, nil
}

func (r *seasonRepository) GetByID(id uint) (*models.Season, error) {
	var season models.Season

	if err := r.db.First(&season, id).Error; err != nil {
		if r.logger != nil {
			r.logger.Error(
				"repository: ошибка получения сезона по id",
				"season_id", id,
				"error", err,
			)
		}
		return nil, err
	}
	return &season, nil
}

func (r *seasonRepository) GetActive() (*models.Season, error) {
	var season models.Season

	if err := r.db.Where("is_active = ?", true).First(&season).Error; err != nil {
		if r.logger != nil {
			r.logger.Error(
				"repository: активный сезон не найден",
				"error", err,
			)
		}
		return nil, err
	}
	return &season, nil
}

func (r *seasonRepository) CloseSeason(id uint) error {
	var season models.Season

	if err := r.db.First(&season, id).Error; err != nil {
		if r.logger != nil {
			r.logger.Error(
				"repository: ошибка закрытия сезона — сезон не найден",
				"season_id", id,
				"error", err,
			)
		}
		return err
	}

	season.IsActive = false

	if err := r.db.Save(&season).Error; err != nil {
		if r.logger != nil {
			r.logger.Error(
				"repository: ошибка сохранения закрытого сезона",
				"season_id", id,
				"error", err,
			)
		}
		return err
	}
	return nil
}

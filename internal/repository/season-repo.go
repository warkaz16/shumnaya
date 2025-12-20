package repository

import (
	"log/slog"
	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type SeasonRepo interface {
	Create(*models.Season) error
	Get() ([]models.Season, error)
	GetByID(ID uint) (*models.Season, error)
	Delete(ID uint) error
	Update(ID uint) error
}

type seasonRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewSeasonRepo(db *gorm.DB, logger *slog.Logger) SeasonRepo {
	return &seasonRepo{
		db:     db,
		logger: logger,
	}
}

func (r *seasonRepo) Create(season *models.Season) error {
	if err := r.db.Create(season).Error; err != nil {
		r.logger.Error(
			"repository.season.Create: ошибка при создании сезона",
			"season_name", season.Name,
			"error", err,
		)
		return err
	}
	return nil
}

func (r *seasonRepo) Get() ([]models.Season, error) {
	var seasons []models.Season
	if err := r.db.Find(&seasons).Error; err != nil {
		r.logger.Error(
			"repository.season.Get: ошибка при получении списка сезонов",
			"error", err,
		)
		return nil, err
	}
	return seasons, nil
}

func (r *seasonRepo) GetByID(ID uint) (*models.Season, error) {
	var season models.Season
	if err := r.db.First(&season, ID).Error; err != nil {
		r.logger.Error(
			"repository.season.GetByID: ошибка при получении сезона",
			"season_id", ID,
			"error", err,
		)
		return nil, err
	}
	return &season, nil
}

func (r *seasonRepo) Update(ID uint) error {
	if err := r.db.Model(&models.Season{}).
		Where("id = ?", ID).
		Updates(map[string]interface{}{
			"is_active": true,
		}).Error; err != nil {
		r.logger.Error(
			"repository.season.Update: ошибка при обновлении сезона",
			"season_id", ID,
			"updated_field", "is_active",
			"error", err,
		)
		return err
	}
	return nil
}

func (r *seasonRepo) Delete(ID uint) error {
	if err := r.db.Delete(&models.Season{}, ID).Error; err != nil {
		r.logger.Error(
			"repository.season.Delete: ошибка при удалении сезона",
			"season_id", ID,
			"error", err,
		)
		return err

	}
	return nil
}

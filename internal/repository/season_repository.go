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
		r.logger.Error("репозиторий: ошибка создания сезона", "error", err)
		return err
	}
	return nil
}

func (r *seasonRepository) GetAll() ([]models.Season, error) {
	var seasons []models.Season
	if err := r.db.Find(&seasons).Error; err != nil {
		r.logger.Error("репозиторий: ошибка получения списка сезонов", "error", err)
		return nil, err
	}
	return seasons, nil
}

func (r *seasonRepository) GetByID(id uint) (*models.Season, error) {
	var season models.Season
	if err := r.db.First(&season, id).Error; err != nil {
		r.logger.Error(
			"репозиторий: ошибка получения сезона по id",
			"season_id", id,
			"error", err,
		)
		return nil, err
	}
	return &season, nil
}

func (r *seasonRepository) GetActive() (*models.Season, error) {
	var season models.Season
	if err := r.db.Where("is_active = ?", true).First(&season).Error; err != nil {
		r.logger.Error("репозиторий: ошибка получения активного сезона", "error", err)
		return nil, err
	}
	return &season, nil
}

func (r *seasonRepository) CloseSeason(id uint) error {
	var season models.Season
	if err := r.db.First(&season, id).Error; err != nil {
		r.logger.Error(
			"репозиторий: ошибка поиска сезона для закрытия",
			"season_id", id,
			"error", err,
		)
		return err
	}

	season.IsActive = false

	if err := r.db.Save(&season).Error; err != nil {
		r.logger.Error(
			"репозиторий: ошибка закрытия сезона",
			"season_id", id,
			"error", err,
		)
		return err
	}

	return nil
}

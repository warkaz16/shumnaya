package repository

import (
	"log/slog"
	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(match *models.Match) error

	Get() ([]models.Match, error)
	GetByID(id uint) (*models.Match, error)

	GetBySeasonID(seasonID uint) ([]models.Match, error)
	GetByPlayerID(playerID uint) ([]models.Match, error)

	GetFiltered(filter *models.MatchFilter) ([]models.Match, error)
}

type matchRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewMatchRepository(db *gorm.DB, log *slog.Logger) MatchRepository {
	return &matchRepository{db: db, log: log}
}

func (r *matchRepository) GetFiltered(filter *models.MatchFilter) ([]models.Match, error) {
	var matches []models.Match

	query := r.db.Model(&matches)

	if filter.SeasonID != nil {
		query = query.Where("season_id = ?", *filter.SeasonID)
	}

	if filter.FromDate != nil {
		query = query.Where("played_at >= ?", *filter.FromDate)
	}

	if filter.ToDate != nil {
		query = query.Where("played_at <= ?", *filter.ToDate)
	}

	err := query.Find(&matches).Error
	return matches, err
}

func (r *matchRepository) Get() ([]models.Match, error) {
	var matches []models.Match
	err := r.db.Find(&matches).Error
	return matches, err
}

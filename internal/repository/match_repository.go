package repository

import (
	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(match *models.Match) error

	GetByID(id uint) (*models.Match, error)

	GetBySeasonID(seasonID uint) ([]models.Match, error)
	GetByPlayerID(playerID uint) ([]models.Match, error)
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) Create(match *models.Match) error {
	return r.db.Create(match).Error
}

func (r *matchRepository) GetByID(id uint) (*models.Match, error) {
	var m models.Match
	if err := r.db.Preload("Winner").Preload("Loser").Preload("Season").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *matchRepository) GetBySeasonID(seasonID uint) ([]models.Match, error) {
	var matches []models.Match
	if err := r.db.Where("season_id = ?", seasonID).
		Preload("Winner").Preload("Loser").Preload("Season").
		Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *matchRepository) GetByPlayerID(playerID uint) ([]models.Match, error) {
	var matches []models.Match
	if err := r.db.Where("winner_id = ? OR loser_id = ?", playerID, playerID).
		Preload("Winner").Preload("Loser").Preload("Season").
		Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

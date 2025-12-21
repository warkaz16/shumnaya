package repository

import (
	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type StandingRepository interface {
	Create(standing *models.Standing) error
	Update(standing *models.Standing) error

	CreateOrUpdate(standing *models.Standing) error

	GetByPlayerAndSeason(playerID, seasonID uint) (*models.Standing, error)
	GetBySeason(seasonID uint) ([]models.Standing, error)
}

type standingRepository struct {
	db *gorm.DB
}

func NewStandingRepository(db *gorm.DB) StandingRepository {
	return &standingRepository{db: db}
}

func (r *standingRepository) Create(standing *models.Standing) error {
	return r.db.Create(standing).Error
}

func (r *standingRepository) Update(standing *models.Standing) error {
	return r.db.Save(standing).Error
}

func (r *standingRepository) CreateOrUpdate(standing *models.Standing) error {
	var existing models.Standing
	err := r.db.Where("player_id = ? AND season_id = ?", standing.PlayerID, standing.SeasonID).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.Create(standing)
		}
		return err
	}

	existing.Wins = standing.Wins
	existing.Losses = standing.Losses
	existing.Points = standing.Points
	existing.Rank = standing.Rank

	return r.Update(&existing)
}

func (r *standingRepository) GetByPlayerAndSeason(playerID, seasonID uint) (*models.Standing, error) {
	var s models.Standing
	if err := r.db.Preload("Player").Preload("Season").Where("player_id = ? AND season_id = ?", playerID, seasonID).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *standingRepository) GetBySeason(seasonID uint) ([]models.Standing, error) {
	var standings []models.Standing
	if err := r.db.Where("season_id = ?", seasonID).Preload("Player").Preload("Season").Find(&standings).Error; err != nil {
		return nil, err
	}
	return standings, nil
}

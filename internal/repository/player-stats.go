package repository

import (
	"context"
	"log/slog"
	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type PlayerStatsRepo interface {
	GetPlayerID(ctx context.Context, userID uint) (*models.PlayerStats, error)
	Update(ctx context.Context, stats *models.PlayerStats) error
}

type playerStatsRepo struct {
	db     *gorm.DB
	logger slog.Logger
}

func NewPlayerStatsRepo(db *gorm.DB, logger slog.Logger) PlayerStatsRepo {
	return &playerStatsRepo{db: db, logger: logger}
}

func (r *playerStatsRepo) GetPlayerID(ctx context.Context, userID uint,
) (*models.PlayerStats, error) {

	var stats models.PlayerStats

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, err
}

func (r *playerStatsRepo) Update(ctx context.Context, stats *models.PlayerStats) error {
	return r.db.WithContext(ctx).Save(&stats).Error
}

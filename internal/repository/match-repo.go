package repository

import (
	"context"
	"fmt"
	"log/slog"
	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type MatchRepo interface {
	GetBySeasonID(ctx context.Context, seasonID uint) ([]models.Match, error)
	GetByID(ctx context.Context, id uint) (*models.Match, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, match *models.Match) error
	Create(ctx context.Context, match *models.Match) error
}

type matchRepo struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewMatchRepo(db *gorm.DB, log *slog.Logger) MatchRepo {
	return &matchRepo{db: db, log: log}
}

func (r *matchRepo) GetBySeasonID(ctx context.Context, seasonID uint) ([]models.Match, error) {
	var matches []models.Match

	err := r.db.WithContext(ctx).Preload("Season").Preload("Player1").Preload("Player2").Preload("Sets").
		Where("season_id = ?", seasonID).Find(&matches).Error

	if err != nil {
		r.log.Error("Ошибка при поиске матчей по seasonId",
			"err", err.Error(),
			"seasonID", seasonID)
		return nil, fmt.Errorf("ошибка при поиске матчей по seasonId: %w", err)
	}

	return matches, nil
}

func (r *matchRepo) GetByID(ctx context.Context, id uint) (*models.Match, error) {
	var match models.Match

	err := r.db.WithContext(ctx).Preload("Season").Preload("Player1").Preload("Player2").Preload("Sets").
		First(&match, id).Error

	if err != nil {
		r.log.Error("Ошибка при поиске матча по Id",
			"err", err.Error(),
			"id", id)
		return nil, fmt.Errorf("ошибка при поиске матча по Id: %w", err)
	}

	return &match, nil
}

func (r *matchRepo) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&models.Match{}, id).Error
	if err != nil {
		r.log.Error("Ошибка при удалении матча",
			"err", err.Error(),
			"id", id)
		return fmt.Errorf("ошибка при удалении матча: %w", err)
	}

	return nil
}

func (r *matchRepo) Update(ctx context.Context, match *models.Match) error {
	err := r.db.WithContext(ctx).Save(match).Error
	if err != nil {
		r.log.Error("Ошибка при обновлении матча",
			"err", err.Error(),
			"matchID", match.ID)
		return fmt.Errorf("ошибка при обновлении матча: %w", err)
	}

	return nil
}

func (r *matchRepo) Create(ctx context.Context, match *models.Match) error {
	err := r.db.WithContext(ctx).Create(match).Error
	if err != nil {
		r.log.Error("Ошибка при создании матча",
			"err", err.Error())
		return fmt.Errorf("ошибка при создании матча: %w", err)
	}

	return nil
}

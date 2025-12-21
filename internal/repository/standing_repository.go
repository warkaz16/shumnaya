package repository

import (
	"log/slog"
	"sort"

	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type StandingRepository interface {
	Create(standing *models.Standing) error
	Update(standing *models.Standing) error

	CreateOrUpdate(standing *models.Standing) error

	GetByPlayerAndSeason(playerID, seasonID uint) (*models.Standing, error)
	GetBySeason(seasonID uint) ([]models.Standing, error)

	// GetSeasonStandingsOrdered возвращает standings отсортированные:
	// 1. По очкам (Points DESC)
	// 2. По разнице побед/поражений (Wins - Losses DESC)
	// 3. По рейтингу игрока (Player.Rating DESC)
	GetSeasonStandingsOrdered(seasonID uint) ([]models.Standing, error)
}

type standingRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewStandingRepository(db *gorm.DB, logger *slog.Logger) StandingRepository {
	return &standingRepository{db: db, logger: logger}
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

func (r *standingRepository) GetSeasonStandingsOrdered(seasonID uint) ([]models.Standing, error) {
	if r.logger != nil {
		r.logger.Info("получение отсортированных standings", "season_id", seasonID)
	}

	var standings []models.Standing

	// Загружаем standings с Player для доступа к рейтингу
	err := r.db.
		Preload("Player").
		Where("season_id = ?", seasonID).
		Find(&standings).Error

	if err != nil {
		if r.logger != nil {
			r.logger.Error("ошибка при получении standings", "season_id", seasonID, "error", err)
		}
		return nil, err
	}

	if r.logger != nil {
		r.logger.Info("standings загружены из БД", "season_id", seasonID, "count", len(standings))
	}

	// Сортируем в Go по критериям:
	// 1. Points DESC
	// 2. (Wins - Losses) DESC
	// 3. Player.Rating DESC
	sort.Slice(standings, func(i, j int) bool {
		// Сначала по очкам
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}

		// Затем по разнице побед/поражений
		diffI := standings[i].Wins - standings[i].Losses
		diffJ := standings[j].Wins - standings[j].Losses
		if diffI != diffJ {
			return diffI > diffJ
		}

		// Затем по рейтингу игрока
		return standings[i].Player.Rating > standings[j].Player.Rating
	})

	if r.logger != nil {
		r.logger.Info("standings отсортированы", "season_id", seasonID, "count", len(standings))
	}

	return standings, nil
}

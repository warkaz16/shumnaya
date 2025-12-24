package repository

import (
	"log/slog"
	"shumnaya/internal/models"

	"gorm.io/gorm"
)

type MatchRepository interface {
	WithDB(tx *gorm.DB) MatchRepository
	Create(match *models.Match) error

	Get() ([]models.Match, error)
	GetByID(id uint) (*models.Match, error)

	GetBySeasonID(seasonID uint) ([]models.Match, error)
	GetByPlayerID(playerID uint) ([]models.Match, error)
	GetRecentByPlayerID(playerID uint, limit int) ([]models.Match, error)

	GetFiltered(filter *models.MatchFilter) ([]models.Match, error)
	HeadToHeadRecordMatchesCount(playerAID, playerBID uint) (int64, error)
	HeadToHeadWinsCount(winnerID, loserID uint) (int64, error)
	HeadToHeadRecentMatches(playerAID, playerBID uint, limit int) ([]models.Match, error)
}

type matchRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewMatchRepository(db *gorm.DB, log *slog.Logger) MatchRepository {
	return &matchRepository{db: db, log: log}
}

func (r *matchRepository) WithDB(tx *gorm.DB) MatchRepository {
	return &matchRepository{db: tx, log: r.log}
}

func (r *matchRepository) GetFiltered(filter *models.MatchFilter) ([]models.Match, error) {
	var matches []models.Match

	query := r.db.Model(&matches)

	if filter.SeasonID != nil {
		query = query.Where("season_id = ?", *filter.SeasonID)
	}

	if filter.PlayerID != nil {
		query = query.Where("winner_id = ? OR loser_id = ?", *filter.PlayerID, *filter.PlayerID)
	}

	if filter.FromDate != nil {
		query = query.Where("played_at >= ?", *filter.FromDate)
	}

	if filter.ToDate != nil {
		query = query.Where("played_at <= ?", *filter.ToDate)
	}

	err := query.Preload("Winner").Preload("Loser").Preload("Season").Find(&matches).Error
	return matches, err
}

func (r *matchRepository) Get() ([]models.Match, error) {
	var matches []models.Match
	// Подгружаем связанные сущности, чтобы в возвращаемых матчах были заполнены
	// поля Winner, Loser и Season
	err := r.db.Preload("Winner").Preload("Loser").Preload("Season").
		Order("played_at DESC").
		Find(&matches).Error
	return matches, err
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

func (r *matchRepository) GetRecentByPlayerID(playerID uint, limit int) ([]models.Match, error) {
	var matches []models.Match
	if err := r.db.Where("winner_id = ? OR loser_id = ?", playerID, playerID).
		Preload("Winner").Preload("Loser").Preload("Season").
		Order("played_at DESC").
		Limit(limit).
		Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *matchRepository) HeadToHeadRecordMatchesCount(playerAID, playerBID uint) (int64, error) {
	var count int64

	if err := r.db.Model(&models.Match{}).
		Where("(winner_id = ? AND loser_id = ?) OR (winner_id = ? AND loser_id = ?)", playerAID, playerBID, playerBID, playerAID).
		Count(&count).Error; err != nil {
		r.log.Error("ошибка получения записи матча между игроками", "error", err)
		return 0, err
	}

	return count, nil
}

func (r *matchRepository) HeadToHeadWinsCount(winnerID, loserID uint) (int64, error) {
	var count int64

	if err := r.db.Model(&models.Match{}).
		Where("winner_id = ? AND loser_id = ?", winnerID, loserID).
		Count(&count).Error; err != nil {
		r.log.Error("ошибка получения количества побед между игроками", "error", err)
		return 0, err
	}

	return count, nil
}

func (r *matchRepository) HeadToHeadRecentMatches(playerAID, playerBID uint, limit int) ([]models.Match, error) {
	var matches []models.Match
	if err := r.db.Model(&models.Match{}).
		Where("(winner_id = ? AND loser_id = ?) OR (winner_id = ? AND loser_id = ?)", playerAID, playerBID, playerBID, playerAID).
		Preload("Winner").Preload("Loser").Preload("Season").Order("played_at DESC").Limit(limit).Find(&matches).Error; err != nil {
		r.log.Error("ошибка получения последних матчей между игроками", "error", err)
		return nil, err
	}

	return matches, nil
}

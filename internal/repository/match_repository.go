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

	GetFiltered(filter *models.MatchFilter) (*models.PaginatedMatches, error)
	HeadToHeadRecordMatchesCount(playerAID, playerBID uint) (countA int64, countB int64, countC int64, err error)

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

func (r *matchRepository) GetFiltered(filter *models.MatchFilter) (*models.PaginatedMatches, error) {

	var matches []models.Match

	page := 1
	pageSize := 50

	if filter.Page != nil && *filter.Page > 0 {
		page = *filter.Page
	}

	if filter.PageSize != nil && *filter.PageSize > 0 && *filter.PageSize <= 1000 {
		pageSize = *filter.PageSize
	}
	
	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Match{})

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

	var total int64
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	err := query.
		Order("played_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&matches).Error
	if err != nil {
		return nil, err
	}

	hasNext := int64(offset+pageSize) < total
    
    return &models.PaginatedMatches{
        Data: matches,
        Pagination: models.Pagination{
            Page:     page,
            PageSize: pageSize,
            Total:    total,
            HasNext:  hasNext,
        },
    }, nil
}

func (r *matchRepository) Get() ([]models.Match, error) {
	var matches []models.Match

	err := r.db.
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

		Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *matchRepository) GetByPlayerID(playerID uint) ([]models.Match, error) {
	var matches []models.Match
	if err := r.db.Where("winner_id = ? OR loser_id = ?", playerID, playerID).

		Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *matchRepository) GetRecentByPlayerID(playerID uint, limit int) ([]models.Match, error) {
	var matches []models.Match
	if err := r.db.Where("winner_id = ? OR loser_id = ?", playerID, playerID).

		Order("played_at DESC").
		Limit(limit).
		Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *matchRepository) HeadToHeadRecordMatchesCount(playerAID, playerBID uint) (countA int64, countB int64, countC int64, err error) {
	var count int64

	if err := r.db.Model(&models.Match{}).
		Where("(winner_id = ? AND loser_id = ?) OR (winner_id = ? AND loser_id = ?)", playerAID, playerBID, playerBID, playerAID).
		Count(&count).Error; err != nil {
		r.log.Error("ошибка получения записи матча между игроками", "error", err)
		return 0, 0, 0, err
	}

	var count1 int64

	if err := r.db.Model(&models.Match{}).
		Where("winner_id = ? AND loser_id = ?", playerAID, playerBID).
		Count(&count1).Error; err != nil {
		r.log.Error("ошибка получения количества побед между игроками", "error", err)
		return 0, 0, 0, err
	}

	var count2 int64

	if err := r.db.Model(&models.Match{}).
		Where("winner_id = ? AND loser_id = ?", playerBID, playerAID).
		Count(&count2).Error; err != nil {
		r.log.Error("ошибка получения количества побед между игроками", "error", err)
		return 0, 0, 0, err
	}

	return count, count1, count2, nil
}

func (r *matchRepository) HeadToHeadRecentMatches(playerAID, playerBID uint, limit int) ([]models.Match, error) {
	var matches []models.Match
	if err := r.db.Model(&models.Match{}).
		Where("(winner_id = ? AND loser_id = ?) OR (winner_id = ? AND loser_id = ?)", playerAID, playerBID, playerBID, playerAID).
		Preload("Winner").Preload("Loser").Order("played_at DESC").Limit(limit).Find(&matches).Error; err != nil {
		r.log.Error("ошибка получения последних матчей между игроками", "error", err)
		return nil, err
	}

	return matches, nil
}

package service

import (
	"errors"
	"time"

	"shumnaya/internal/models"
	"shumnaya/internal/repository"
	"shumnaya/internal/utils/elo"

	"gorm.io/gorm"
)

type MatchService interface {
	RecordMatch(winnerID, loserID, seasonID uint, score string) (*models.Match, error)
}

type matchService struct {
	db           *gorm.DB
	matchRepo    repository.MatchRepository
	playerRepo   repository.PlayerRepository
	standingRepo repository.StandingRepository
}

func NewMatchService(db *gorm.DB, mr repository.MatchRepository, pr repository.PlayerRepository, sr repository.StandingRepository) MatchService {
	return &matchService{db: db, matchRepo: mr, playerRepo: pr, standingRepo: sr}
}

func (s *matchService) RecordMatch(winnerID, loserID, seasonID uint, score string) (*models.Match, error) {
	if winnerID == loserID {
		return nil, errors.New("winner and loser cannot be the same")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	matchRepoTx := repository.NewMatchRepository(tx)
	standingRepoTx := repository.NewStandingRepository(tx)

	var winner models.Player
	if err := tx.First(&winner, winnerID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var loser models.Player
	if err := tx.First(&loser, loserID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	wStanding, err := standingRepoTx.GetByPlayerAndSeason(winnerID, seasonID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, err
	}
	if wStanding == nil {
		wStanding = &models.Standing{PlayerID: winnerID, SeasonID: seasonID}
	}

	lStanding, err := standingRepoTx.GetByPlayerAndSeason(loserID, seasonID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, err
	}
	if lStanding == nil {
		lStanding = &models.Standing{PlayerID: loserID, SeasonID: seasonID}
	}

	wGames := wStanding.Wins + wStanding.Losses
	lGames := lStanding.Wins + lStanding.Losses

	winnerNew := elo.NewRating(winner.Rating, loser.Rating, true, wGames, wStanding.Wins)
	loserNew := elo.NewRating(loser.Rating, winner.Rating, false, lGames, lStanding.Wins)

	winnerChange := winnerNew - winner.Rating
	loserChange := loserNew - loser.Rating

	match := &models.Match{
		WinnerID:           winnerID,
		LoserID:            loserID,
		SeasonID:           seasonID,
		Score:              score,
		WinnerRatingChange: winnerChange,
		LoserRatingChange:  loserChange,
		PlayedAt:           time.Now(),
	}

	if err := matchRepoTx.Create(match); err != nil {
		tx.Rollback()
		return nil, err
	}

	winner.Rating = winnerNew
	loser.Rating = loserNew

	if err := tx.Save(&winner).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Save(&loser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	wStanding.Wins += 1
	wStanding.Points += 1

	lStanding.Losses += 1

	if err := standingRepoTx.CreateOrUpdate(wStanding); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := standingRepoTx.CreateOrUpdate(lStanding); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return match, nil
}

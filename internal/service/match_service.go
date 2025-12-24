package service

import (
	"errors"
	"log/slog"
	"time"

	"shumnaya/internal/models"
	"shumnaya/internal/repository"
	"shumnaya/internal/utils/elo"

	"gorm.io/gorm"
)

type MatchService interface {
	RecordMatch(winnerID, loserID, seasonID uint, score string) (*models.Match, error)

	Get() ([]models.Match, error)
	GetFiltered(filter *models.MatchFilter) ([]models.Match, error)
	GetHeadToHead(playerAID, playerBID uint, limit int) (*models.HeadToHeadRecord, error)
}

type matchService struct {
	db           *gorm.DB
	logger       *slog.Logger
	matchRepo    repository.MatchRepository
	playerRepo   repository.PlayerRepository
	standingRepo repository.StandingRepository
}

func NewMatchService(db *gorm.DB, log *slog.Logger, mr repository.MatchRepository, pr repository.PlayerRepository, sr repository.StandingRepository) MatchService {
	return &matchService{db: db, logger: log, matchRepo: mr, playerRepo: pr, standingRepo: sr}
}

func (s *matchService) RecordMatch(winnerID, loserID, seasonID uint, score string) (*models.Match, error) {
	if winnerID == loserID {
		return nil, errors.New("winner and loser cannot be the same")
	}

	var created *models.Match

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Репозитории в контексте транзакции
		matchRepoTx := s.matchRepo.WithDB(tx)
		standingRepoTx := s.standingRepo.WithDB(tx)

		var winner models.Player
		if err := tx.First(&winner, winnerID).Error; err != nil {
			return err
		}

		var loser models.Player
		if err := tx.First(&loser, loserID).Error; err != nil {
			return err
		}

		wStanding, err := standingRepoTx.GetByPlayerAndSeason(winnerID, seasonID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if wStanding == nil || err == gorm.ErrRecordNotFound {
			wStanding = &models.Standing{PlayerID: winnerID, SeasonID: seasonID}
		}

		lStanding, err := standingRepoTx.GetByPlayerAndSeason(loserID, seasonID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if lStanding == nil || err == gorm.ErrRecordNotFound {
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
			return err
		}

		winner.Rating = winnerNew
		loser.Rating = loserNew

		if err := tx.Save(&winner).Error; err != nil {
			return err
		}
		if err := tx.Save(&loser).Error; err != nil {
			return err
		}

		wStanding.Wins += 1
		wStanding.Points += 1

		lStanding.Losses += 1

		if err := standingRepoTx.CreateOrUpdate(wStanding); err != nil {
			return err
		}
		if err := standingRepoTx.CreateOrUpdate(lStanding); err != nil {
			return err
		}

		created = match
		return nil
	})

	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *matchService) GetFiltered(filter *models.MatchFilter) ([]models.Match, error) {
	return s.matchRepo.GetFiltered(filter)
}

func (s *matchService) Get() ([]models.Match, error) {
	// Логирование начала операции
	s.logger.Info("получение всех матчей")

	// Получение матчей из репозитория
	matches, err := s.matchRepo.Get()
	if err != nil {
		s.logger.Error("ошибка при получении матчей", "ошибка", err)
		return nil, err
	}

	// Валидация результата
	if matches == nil {
		return []models.Match{}, nil
	}

	s.logger.Info("матчи успешно получены", "количество", len(matches))
	return matches, nil
}

func (s *matchService) GetHeadToHead(playerAID, playerBID uint, limit int) (*models.HeadToHeadRecord, error) {
	var record models.HeadToHeadRecord

	record.PlayerAID = playerAID
	record.PlayerBID = playerBID

	totalMatches, err := s.matchRepo.HeadToHeadRecordMatchesCount(playerAID, playerBID)
	if err != nil {
		return nil, err
	}
	record.TotalMatches = int(totalMatches)

	playerAWins, err := s.matchRepo.HeadToHeadWinsCount(playerAID, playerBID)
	if err != nil {
		return nil, err
	}
	record.PlayerAWins = int(playerAWins)

	playerBWins, err := s.matchRepo.HeadToHeadWinsCount(playerBID, playerAID)
	if err != nil {
		return nil, err
	}
	record.PlayerBWins = int(playerBWins)

	recentMatches, err := s.matchRepo.HeadToHeadRecentMatches(playerAID, playerBID, limit)
	if err != nil {
		return nil, err
	}
	record.LastMatchesPlayed = recentMatches

	return &record, nil
}
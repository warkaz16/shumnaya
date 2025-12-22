package service

import (
	"fmt"
	"log/slog"
	"shumnaya/internal/models"
	"shumnaya/internal/repository"
)

type StandingService interface {
	GetSeasonStandings(seasonID uint) ([]models.Standing, error)
}

type standingService struct {
	repo repository.StandingRepository
	log  *slog.Logger
}

func NewStandingService(repo repository.StandingRepository, log *slog.Logger) StandingService {
	return &standingService{repo: repo, log: log}
}

func (s *standingService) GetSeasonStandings(seasonID uint) ([]models.Standing, error) {

	standings, err := s.repo.GetSeasonStandingsOrdered(seasonID)

	if err != nil {
		s.log.Error(
			"service: ошибка при получении турнирной таблицы сезона",
			"season_id", seasonID,
			"error", err)
		return nil, fmt.Errorf("Ошибка при получении Standings: %w", err)
	}

	s.log.Info(
		"service: турнирная таблица успешно получена",
		"season_id", seasonID,
		"count", len(standings),
	)

	return standings, nil

}

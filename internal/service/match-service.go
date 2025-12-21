package service

import (
	"errors"
	"log/slog"
	"shumnaya/internal/models"
	"shumnaya/internal/repository"
)

type MatchService interface {
	Get() ([]models.Match, error)
	GetFiltered(filter *models.MatchFilter) ([]models.Match, error)
}

type matchService struct {
	repo   repository.MatchRepository
	logger *slog.Logger
}

func NewMatchService(repo repository.MatchRepository, logger *slog.Logger) MatchService {
	return &matchService{
		repo:   repo,
		logger: logger,
	}
}

func (s *matchService) GetFiltered(filter *models.MatchFilter) ([]models.Match, error) {
	return s.repo.GetFiltered(filter)
}

func (s *matchService) Get() ([]models.Match, error) {
	// Логирование начала операции
	s.logger.Info("получение всех матчей")

	// Получение матчей из репозитория
	matches, err := s.repo.Get()
	if err != nil {
		s.logger.Error("ошибка при получении матчей", "ошибка", err)
		return nil, err
	}

	// Валидация результата
	if matches == nil {
		errors.New("матчей нет")
	}

	s.logger.Info("матчи успешно получены", "количество", len(matches))
	return matches, nil
}

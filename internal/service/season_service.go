package service

import (
	"errors"
	"log/slog"

	"shumnaya/internal/models"
	"shumnaya/internal/repository"
)

type SeasonService interface {
	CreateSeason(season *models.Season) error
	GetAllSeasons() ([]models.Season, error)
	GetSeasonByID(id uint) (*models.Season, error)
}

type seasonService struct {
	repo   repository.SeasonRepository
	logger *slog.Logger
}

func NewSeasonService(
	repo repository.SeasonRepository,
	logger *slog.Logger,
) SeasonService {
	return &seasonService{
		repo:   repo,
		logger: logger,
	}
}

func (s *seasonService) CreateSeason(season *models.Season) error {
	if season.StartDate.After(season.EndDate) {
		err := errors.New("дата начала должна быть раньше даты окончания")

		if s.logger != nil {
			s.logger.Error(
				"service: некорректные даты сезона",
				"error", err,
			)
		}

		return err
	}

	season.IsActive = true

	if err := s.repo.Create(season); err != nil {
		if s.logger != nil {
			s.logger.Error(
				"service: ошибка создания сезона",
				"error", err,
			)
		}
		return err
	}

	return nil
}

func (s *seasonService) GetAllSeasons() ([]models.Season, error) {
	seasons, err := s.repo.GetAll()
	if err != nil {
		if s.logger != nil {
			s.logger.Error(
				"service: ошибка получения списка сезонов",
				"error", err,
			)
		}
		return nil, err
	}
	return seasons, nil
}

func (s *seasonService) GetSeasonByID(id uint) (*models.Season, error) {
	season, err := s.repo.GetByID(id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(
				"service: ошибка получения сезона по id",
				"season_id", id,
				"error", err,
			)
		}
		return nil, err
	}
	return season, nil
}

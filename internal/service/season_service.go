package service

import (
	"errors"
	"log/slog"

	"shumnaya/internal/models"
	"shumnaya/internal/repository"
)

type SeasonService interface {
	CreateSeason(season *models.Season) error
	GetAll() ([]models.Season, error)
	GetByID(id uint) (*models.Season, error)
	GetActive() (*models.Season, error)
	CloseSeason(id uint) error
}

type seasonService struct {
	repo   repository.SeasonRepository
	logger *slog.Logger
}

func NewSeasonService(repo repository.SeasonRepository, logger *slog.Logger) SeasonService {
	return &seasonService{
		repo:   repo,
		logger: logger,
	}
}

func (s *seasonService) CreateSeason(season *models.Season) error {
	if season.StartDate.After(season.EndDate) {
		if s.logger != nil {
			s.logger.Info("ошибка валидации сезона: дата начала позже даты окончания")
		}
		return errors.New("дата начала сезона должна быть раньше даты окончания")
	}

	season.IsActive = true

	if s.logger != nil {
		s.logger.Info("сервис: создание сезона", "name", season.Name)
	}

	return s.repo.Create(season)
}

func (s *seasonService) GetAll() ([]models.Season, error) {
	return s.repo.GetAll()
}

func (s *seasonService) GetByID(id uint) (*models.Season, error) {
	return s.repo.GetByID(id)
}

func (s *seasonService) GetActive() (*models.Season, error) {
	return s.repo.GetActive()
}

func (s *seasonService) CloseSeason(id uint) error {
	if s.logger != nil {
		s.logger.Info("сервис: закрытие сезона", "season_id", id)
	}
	return s.repo.CloseSeason(id)
}

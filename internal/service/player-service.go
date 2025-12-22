package service

import (
	"errors"
	"log/slog"
 "golang.org/x/crypto/bcrypt"

	"shumnaya/internal/models"
	"shumnaya/internal/repository"

	"gorm.io/gorm"
)

type PlayerService interface {
	GetPlayerProfile(id uint) (*models.PlayerProfile, error)
		RegisterPlayer(name, email, password string) error

}

const defaultRecentMatchesLimit = 5

type playerService struct {
	db         *gorm.DB
	logger     *slog.Logger
	playerRepo repository.PlayerRepository
	matchRepo  repository.MatchRepository
}

func NewPlayerService(db *gorm.DB, log *slog.Logger, pr repository.PlayerRepository, mr repository.MatchRepository) PlayerService {
	return &playerService{db: db, logger: log, playerRepo: pr, matchRepo: mr}
}

func (s *playerService) RegisterPlayer(
	name string,
	email string,
	password string,
) error {

	if _, err := s.playerRepo.GetByEmail(email); err == nil {
		return errors.New("игрок с таким email уже существует")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(
				"service: ошибка хеширования пароля",
				"error", err,
			)
		}
		return err
	}

	player := &models.Player{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Rating:       1000,
	}

	if err := s.playerRepo.Create(player); err != nil {
		if s.logger != nil {
			s.logger.Error(
				"service: ошибка регистрации игрока",
				"email", email,
				"error", err,
			)
		}
		return err
	}

	return nil
}


func (s *playerService) GetPlayerProfile(id uint) (*models.PlayerProfile, error) {
	if id == 0 {
		return nil, errors.New("invalid player id")
	}

	player, err := s.playerRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service: player not found", "player_id", id, "error", err)
		return nil, err
	}

	matches, err := s.matchRepo.GetByPlayerID(id)
	if err != nil {
		s.logger.Error("service: failed to get player matches", "player_id", id, "error", err)
		return nil, err
	}

	recentMatches, err := s.matchRepo.GetRecentByPlayerID(id, defaultRecentMatchesLimit)
	if err != nil {
		s.logger.Error("service: failed to get recent player matches", "player_id", id, "error", err)
		return nil, err
	}

	total := len(matches)
	wins := 0
	for _, m := range matches {

		if m.WinnerID == id {
			wins++
		}
	}

	losses := total - wins

	profile := &models.PlayerProfile{
		Player:        *player,
		Rating:        player.Rating,
		TotalMatches:  total,
		Wins:          wins,
		Losses:        losses,
		RecentMatches: recentMatches,
	}

	return profile, nil
}

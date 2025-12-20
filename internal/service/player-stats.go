package service

import (
	"context"
	"log/slog"
	"shumnaya/internal/domain/elo"
	"shumnaya/internal/repository"
)

type StatsService interface {
	ApplyMatchResult(ctx context.Context, playerID uint, opponentID uint, isWin bool) error
}

type statsService struct {
	repo   repository.PlayerStatsRepo
	logger slog.Logger
}

func NewPlayerStatsService(repo repository.PlayerStatsRepo, logger slog.Logger) StatsService {
	return &statsService{repo: repo, logger: logger}
}

func (s *statsService) ApplyMatchResult(ctx context.Context, playerID uint, opponentID uint, isWin bool,
) error {

	player, err := s.repo.GetPlayerID(ctx, playerID)
	if err != nil {
		return err
	}

	opponent, err := s.repo.GetPlayerID(ctx, opponentID)
	if err != nil {
		return err
	}

	player.GamesPlayed++
	if isWin {
		player.Wins++
	} else {
		player.Loses++
	}

	player.Elo = elo.NewRating(
		player.Elo,
		opponent.Elo,
		isWin,
		player.GamesPlayed,
		player.Wins,
	)

	return s.repo.Update(ctx, player)
}

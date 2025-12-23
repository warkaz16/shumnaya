package transport

import (
	"log/slog"
	"shumnaya/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	matchService service.MatchService,
	playerService service.PlayerService,
	seasonService service.SeasonService,
	standingService service.StandingService,
	logger *slog.Logger,
) {
	matchHandler := NewMatchHandler(r, matchService, logger)
	playerHandler := NewPlayerHandler(r, playerService, logger)
	seasonHandler := NewSeasonHandler(r, seasonService, standingService, logger)

	matchHandler.RegisterRoutes(r)
	playerHandler.RegisterRoutes(r)
	seasonHandler.RegisterRoutes(r)
}

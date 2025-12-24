package transport

import (
	"log/slog"
	"shumnaya/internal/service"
	"shumnaya/internal/transport/middleware"

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

	// все как было
	matchHandler.RegisterRoutes(r)
	seasonHandler.RegisterRoutes(r)

	// 🔓 публичные
	r.POST("/players", playerHandler.Register)
	r.POST("/login", playerHandler.Login)

	// 🔐 защищённые
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/players/:id", playerHandler.GetByID)
}

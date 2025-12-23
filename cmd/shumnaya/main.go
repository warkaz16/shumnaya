package main

import (
	"log"

	"shumnaya/internal/config"
	"shumnaya/internal/models"
	"shumnaya/internal/repository"
	"shumnaya/internal/service"
	"shumnaya/internal/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	logger := config.InitLogger()

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	db := config.ConnectDB(logger)

	if err := db.AutoMigrate(&models.Match{}, &models.Player{}, &models.Season{}, &models.Standing{}); err != nil {
		logger.Error("ошибка миграции базы данных", "error", err)
		log.Fatal("Ошибка миграции базы данных:", err)
	}

	logger.Info("Миграция базы данных выполнена успешно")

	matchRepo := repository.NewMatchRepository(db, logger)
	seasonRepo := repository.NewSeasonRepository(db, logger)
	playerRepo := repository.NewPlayerRepository(db, logger)
	standingRepo := repository.NewStandingRepository(db, logger)

	matchService := service.NewMatchService(db, logger, matchRepo, playerRepo, standingRepo)
	playerService := service.NewPlayerService(db, logger, playerRepo, matchRepo)
	seasonService := service.NewSeasonService(seasonRepo, logger)
	standingService := service.NewStandingService(standingRepo, logger)

	r := gin.Default()

	transport.RegisterRoutes(
		r, matchService, playerService, seasonService, standingService, logger,
	)

	logger.Info("Server running on :8080")

	if err := r.Run(); err != nil {
		logger.Error("ошибка запуска сервера", "error", err)
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

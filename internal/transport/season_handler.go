package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"shumnaya/internal/models"
	"shumnaya/internal/service"

	"github.com/gin-gonic/gin"
)

type SeasonHandler struct {
	service  service.SeasonService
	standing service.StandingService
	logger   *slog.Logger
}

func NewSeasonHandler(
	r *gin.Engine,
	service service.SeasonService,
	standing service.StandingService,
	logger *slog.Logger,

) *SeasonHandler {
	return &SeasonHandler{
		service:  service,
		logger:   logger,
		standing: standing,
	}
}

func (h *SeasonHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/seasons", h.getAll)
	r.GET("/seasons/:id", h.getByID)
	r.POST("/seasons", h.create)
	r.GET("/seasons/:id/standings", h.getByIDstandings)
}

func (h *SeasonHandler) getAll(c *gin.Context) {
	seasons, err := h.service.GetAllSeasons()
	if err != nil {
		h.logger.Error("handler: ошибка получения списка сезонов", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить сезоны"})
		return
	}

	c.JSON(http.StatusOK, seasons)
}

func (h *SeasonHandler) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		h.logger.Error("handler: некорректный id сезона")
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	season, err := h.service.GetSeasonByID(uint(id))
	if err != nil {
		h.logger.Error("handler: сезон не найден", "season_id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "сезон не найден"})
		return
	}

	c.JSON(http.StatusOK, season)
}

func (h *SeasonHandler) create(c *gin.Context) {
	var season models.Season

	if err := c.ShouldBindJSON(&season); err != nil {
		h.logger.Error("handler: ошибка валидации данных сезона", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateSeason(&season); err != nil {
		h.logger.Error("handler: ошибка создания сезона", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, season)
}

func (h *SeasonHandler) getByIDstandings(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		h.logger.Error("handler: некорректный id сезона")
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	season, err := h.service.GetSeasonByID(uint(id))
	if err != nil {
		h.logger.Error("handler: сезон не найден", "season_id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "сезон не найден"})
		return
	}

	standing, err := h.standing.GetSeasonStandings(season.ID)

	if err != nil {
		h.logger.Error("Ошибка при поиске standing")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(standing) == 0 {
		h.logger.Error("season standings пустые")
		c.JSON(http.StatusOK, gin.H{
			"season_id": id,
			"standings": standing,
			"count":     len(standing),
		})
		return
	}

	c.JSON(http.StatusOK, standing)

}

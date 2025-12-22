package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"shumnaya/internal/models"
	"shumnaya/internal/service"

	"github.com/gin-gonic/gin"
)

type seasonHandler struct {
	service service.SeasonService
	logger  *slog.Logger
}

func NewSeasonHandler(
	r *gin.Engine,
	service service.SeasonService,
	logger *slog.Logger,
) {
	h := &seasonHandler{
		service: service,
		logger:  logger,
	}

	r.GET("/seasons", h.getAll)
	r.GET("/seasons/:id", h.getByID)
	r.POST("/seasons", h.create)
}

func (h *seasonHandler) create(c *gin.Context) {
	var season models.Season

	if err := c.ShouldBindJSON(&season); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректные данные сезона",
		})
		return
	}

	if err := h.service.CreateSeason(&season); err != nil {
		if h.logger != nil {
			h.logger.Error(
				"handler: ошибка создания сезона",
				"error", err,
			)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, season)
}


func (h *seasonHandler) getAll(c *gin.Context) {
	seasons, err := h.service.GetAllSeasons()
	if err != nil {
		if h.logger != nil {
			h.logger.Error(
				"handler: ошибка получения списка сезонов",
				"error", err,
			)
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось получить сезоны",
		})
		return
	}

	c.JSON(http.StatusOK, seasons)
}

func (h *seasonHandler) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный id",
		})
		return
	}

	season, err := h.service.GetSeasonByID(uint(id))
	if err != nil {
		if h.logger != nil {
			h.logger.Error(
				"handler: ошибка получения сезона",
				"season_id", id,
				"error", err,
			)
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error": "сезон не найден",
		})
		return
	}

	c.JSON(http.StatusOK, season)
}



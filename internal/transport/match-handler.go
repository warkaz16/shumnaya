package transport

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"shumnaya/internal/models"
	"shumnaya/internal/service"

	"github.com/gin-gonic/gin"
)

type matchHandler struct {
	service service.MatchService
	logger  *slog.Logger
}

func NewMatchHandler(service service.MatchService, logger *slog.Logger) *matchHandler {
	return &matchHandler{
		service: service,
		logger:  logger,
	}
}

func (h *matchHandler) GetMatches(c *gin.Context) {
	filter := &models.MatchFilter{}

	if seasonIDStr := c.Query("season_id"); seasonIDStr != "" {
		if seasonID, err := strconv.ParseUint(seasonIDStr, 10, 32); err == nil {
			seasonIDUint := uint(seasonID)
			filter.SeasonID = &seasonIDUint
			h.logger.Info("фильтр по сезону", "season_id", seasonIDUint)
		} else {

			h.logger.Warn("некорректный параметр season_id", "значение", seasonIDStr, "ошибка", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid season_id format"})
			return
		}
	}

	if playerIDStr := c.Query("player_id"); playerIDStr != "" {
		if playerID, err := strconv.ParseUint(playerIDStr, 10, 32); err == nil {
			playerIDUint := uint(playerID)
			filter.PlayerID = &playerIDUint
			h.logger.Info("фильтр по игроку", "player_id", playerIDUint)
		} else {
			h.logger.Warn("некорректный параметр player_id", "значение", playerIDStr, "ошибка", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player_id format"})
			return
		}
	}

	if fromStr := c.Query("from"); fromStr != "" {
		if fromTime, err := time.Parse("02.01.06", fromStr); err == nil {
			filter.FromDate = &fromTime
			h.logger.Info("фильтр по начальной дате", "from", fromTime)
		} else {
			h.logger.Warn("некорректный параметр from", "значение", fromStr, "ошибка", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный формат даты from, используйте ДД.МММ.ГГ (например: 25.12.24)"})
			return
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if toTime, err := time.Parse("02.01.06", toStr); err == nil {
			filter.ToDate = &toTime
			h.logger.Info("фильтр по конечной дате", "to", toTime)
		} else {
			h.logger.Warn("некорректный параметр to", "значение", toStr, "ошибка", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный формат даты to, используйте ДД.ММ.ГГ (например: 25.12.24)"})
			return
		}
	}

	matches, err := h.service.GetFiltered(filter)
	if err != nil {
		h.logger.Error("ошибка при получении матчей", "ошибка", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch matches"})
		return
	}

	h.logger.Info("матчи успешно получены")

	c.JSON(http.StatusOK, gin.H{
		"data": matches,
	})
}

func (h *matchHandler) CreateMatch(c *gin.Context) {
	var req models.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.WinnerID == req.LoserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "winner_id and loser_id must be different"})
		return
	}

	match, err := h.service.RecordMatch(req.WinnerID, req.LoserID, req.SeasonID, req.Score)
	if err != nil {
		h.logger.Error("failed to record match", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record match"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": match})
}

package transport

import (
	"log/slog"
	"strconv"

	"shumnaya/internal/service"

	"github.com/gin-gonic/gin"
)

type playerHandler struct {
	service service.PlayerService
	logger  *slog.Logger
}

type registerPlayerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewPlayerHandler(r *gin.Engine, svc service.PlayerService, logger *slog.Logger) {
	h := &playerHandler{
		service: svc,
		logger:  logger,
	}

	r.GET("/players/:id", h.getByID)
	r.POST("/players", h.register)
}

func (h *playerHandler) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{
			"error": "некорректный id игрока",
		})
		return
	}

	profile, err := h.service.GetPlayerProfile(uint(id))
	if err != nil {
		if h.logger != nil {
			h.logger.Error(
				"handler: ошибка получения профиля игрока",
				"player_id", id,
				"error", err,
			)
		}
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, profile)
}

func (h *playerHandler) register(c *gin.Context) {
	var req registerPlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "некорректные данные регистрации",
		})
		return
	}

	if err := h.service.RegisterPlayer(
		req.Name,
		req.Email,
		req.Password,
	); err != nil {
		if h.logger != nil {
			h.logger.Error(
				"handler: ошибка регистрации игрока",
				"email", req.Email,
				"error", err,
			)
		}
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "игрок успешно зарегистрирован",
	})
}

package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"shumnaya/internal/service"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	service service.PlayerService
	logger  *slog.Logger
}

func NewPlayerHandler(r *gin.Engine, svc service.PlayerService, logger *slog.Logger) *PlayerHandler {
	return &PlayerHandler{service: svc, logger: logger}
}

func (h *PlayerHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/players/:id", h.GetByID)
	r.POST("/players", h.Register)
}

// GetByID godoc
// @Summary Профиль игрока
// @Tags Players
// @Produce json
// @Param id path int true "ID игрока"
// @Success 200 {object} models.PlayerProfile
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /players/{id} [get]
func (h *PlayerHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn("handler: invalid player id", "id", idStr, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	profile, err := h.service.GetPlayerProfile(uint(id))
	if err != nil {
		h.logger.Error("handler: failed to get player profile", "player_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

type RegisterPlayerRequest struct {
	Name     string `json:"name" binding:"required" example:"Иван"`
	Email    string `json:"email" binding:"required,email" example:"test@mail.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
}

type RegisterPlayerResponse struct {
	Message string `json:"message" example:"игрок успешно зарегистрирован"`
}

// Register godoc
// @Summary Регистрация игрока
// @Description Создает нового игрока
// @Tags Players
// @Accept json
// @Produce json
// @Param input body RegisterPlayerRequest true "Данные регистрации"
// @Success 201 {object} RegisterPlayerResponse
// @Failure 400 {object} map[string]string
// @Router /players [post]
func (h *PlayerHandler) Register(c *gin.Context) {
	var req RegisterPlayerRequest

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

package transport

import (
	"log/slog"
	"strconv"

	"shumnaya/internal/dto"
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
	tokenPlayerID := c.GetUint("player_id")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	if uint(id) != tokenPlayerID {
		c.JSON(403, gin.H{"error": "доступ запрещён"})
		return
	}

	profile, err := h.service.GetPlayerProfile(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, profile)
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
	var req dto.RegisterPlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "некорректные данные регистрации",
		})
		return
	}

	token, err := h.service.RegisterPlayer(
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {
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
		"token": token,
	})
}

func (h *PlayerHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "некорректные данные"})
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.LoginResponse{
		Token: token,
	})
}

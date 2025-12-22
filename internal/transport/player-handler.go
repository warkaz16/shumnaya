package transport

import (
    "log/slog"
    "net/http"
    "strconv"

    "shumnaya/internal/service"

    "github.com/gin-gonic/gin"
)

type playerHandler struct {
    service service.PlayerService
    logger  *slog.Logger
}

func NewPlayerHandler(r *gin.Engine, svc service.PlayerService, logger *slog.Logger) {
    h := &playerHandler{service: svc, logger: logger}

    r.GET("/players/:id", h.getByID)
}

func (h *playerHandler) getByID(c *gin.Context) {
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

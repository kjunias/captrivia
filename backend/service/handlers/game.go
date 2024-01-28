package handlers

import (
	"net/http"

	"github.com/ProlificLabs/captrivia/backend/service"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	router     *gin.Engine
	gameServer *service.GameServer
}

func NewGameHandler(r *gin.Engine, gs *service.GameServer) *GameHandler {
	return &GameHandler{
		router:     r,
		gameServer: gs,
	}
}

func (h *GameHandler) RegisterRoutes() {
	h.router.POST("/game/start", h.handleStart)
	h.router.POST("/game/end", h.handleEnd)
}

func (h *GameHandler) handleStart(c *gin.Context) {
	sessionID := h.gameServer.Sessions.CreateSession()
	c.JSON(http.StatusOK, gin.H{"sessionId": sessionID})
}

func (h *GameHandler) handleEnd(c *gin.Context) {
	var request struct {
		SessionID string `json:"sessionId"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, exists := h.gameServer.Sessions.GetSession(request.SessionID)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"finalScore": session.Score})
}

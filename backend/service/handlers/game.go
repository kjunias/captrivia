package handlers

import (
	"fmt"
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
	h.router.POST("/gameroom/create", h.handleCreateGameRoom)
	h.router.GET("/gameroom/join", h.handleJoinGameRoom)
	h.router.GET("/gameroom/update", h.handleUpdateGameRoom)
	h.router.POST("/gameroom/start", h.handleStartGameRoom)
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

func (h *GameHandler) handleCreateGameRoom(c *gin.Context) {
	h.gameServer.CreateGameRoom(c)
}

func (h *GameHandler) handleJoinGameRoom(c *gin.Context) {
	roomID, ok := c.GetQuery("roomID")
	if !ok {
		c.JSON(http.StatusNotFound, "Could not find game room")
		return
	}
	fmt.Println("==> GameHandler handleJoinGameRoom 0: ", roomID)
	h.gameServer.JoinGameRoom(roomID, c, false)
	fmt.Println("==> GameHandler handleJoinGameRoom 1: ", roomID)
}

func (h *GameHandler) handleUpdateGameRoom(c *gin.Context) {
	roomID, ok := c.GetQuery("roomID")
	playerID, ok := c.GetQuery("playerID")
	if !ok {
		fmt.Println("==> handleUpdateGameRoom Error: ")
		c.JSON(http.StatusNotFound, "Error handling update request")
		return
	}
	fmt.Println("==> handleUpdateGameRoom 0:")
	h.gameServer.SubscribeToGameRoom(roomID, playerID, c)
}

func (h *GameHandler) handleStartGameRoom(c *gin.Context) {
	var startInfo struct {
		RoomID            string `json:"roomID"`
		NumberOfQuestions int    `json:"numberOfQuestions"`
	}
	fmt.Println("=====> handleStartGameRoom start")
	if err := c.ShouldBindJSON(&startInfo); err != nil {
		fmt.Println("=====> handleStartGameRoom start error 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	if startInfo.NumberOfQuestions < 1 {
		fmt.Println("=====> handleStartGameRoom start error 1")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of questions"})
		return
	}

	h.gameServer.StartGameRoom(startInfo.RoomID, startInfo.NumberOfQuestions, c)

}

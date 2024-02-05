package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ProlificLabs/captrivia/backend/service"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type GameHandler struct {
	router     *gin.Engine
	gameServer *service.GameServer
}

type IncomingMessage struct {
	Action            string `json:"action"`
	RoomID            string `json:"roomID"`
	PlayerID          string `json:"playerID"`
	QuestionId        string `json:"questionId"`
	Answer            int    `json:"answer"`
	NumberOfQuestions int    `json:"numberOfQuestions"`
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
	h.router.GET("/gameroom/websocket", h.handleWS)
	h.router.POST("/gameroom/start", h.handleStartGameRoom)

	hub := h.gameServer.Hub

	hub.HandleMessage(func(s *melody.Session, msg []byte) {
		var im IncomingMessage
		err := json.Unmarshal(msg, &im)
		if err != nil {
			fmt.Println("Error unmarshalling message:", msg, err)
			return
		}
		roomID := im.RoomID
		switch im.Action {
		case "START_COUNTER":
			h.handleStartCounter(roomID, im)
		case "END_COUNTER":
			h.handleEndCounter(roomID)
		case "SUBMIT_ANSWER":
			h.handleSubmitAnswer(roomID, im)
		default:
			fmt.Println("Unknown actoin")
		}
	})
}

func (h *GameHandler) handleWS(c *gin.Context) {
	h.gameServer.Hub.HandleRequest(c.Writer, c.Request)
}

func (h *GameHandler) handleStartCounter(roomID string, msg IncomingMessage) {
	h.gameServer.HandleStartCounter(roomID, msg.NumberOfQuestions)
}

func (h *GameHandler) handleEndCounter(roomID string) {
	h.gameServer.HandleEndCounter(roomID)
}

func (h *GameHandler) handleSubmitAnswer(roomID string, msg IncomingMessage) {
	h.gameServer.HandleSubmitAnswer(roomID, msg.PlayerID, msg.QuestionId, msg.Answer)
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
	h.gameServer.JoinGameRoom(roomID, c, false)
}

func (h *GameHandler) handleStartGameRoom(c *gin.Context) {
	var startInfo struct {
		RoomID            string `json:"roomID"`
		NumberOfQuestions int    `json:"numberOfQuestions"`
	}
	if err := c.ShouldBindJSON(&startInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	if startInfo.NumberOfQuestions < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of questions"})
		return
	}

	h.gameServer.StartGameRoom(startInfo.RoomID, startInfo.NumberOfQuestions, c)
}

package service

import (
	"net/http"

	"github.com/ProlificLabs/captrivia/backend/model"
	"github.com/ProlificLabs/captrivia/backend/session"
	"github.com/gin-gonic/gin"
)

type GameServer struct {
	Questions []model.Question
	Sessions  *session.SessionStore
}

func NewGameServer(questions []model.Question, store *session.SessionStore) *GameServer {
	return &GameServer{
		Questions: questions,
		Sessions:  store,
	}
}

func (gs *GameServer) EndGameHandler(c *gin.Context) {
	var request struct {
		SessionID string `json:"sessionId"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, exists := gs.Sessions.GetSession(request.SessionID)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"finalScore": session.Score})
}

func (gs *GameServer) StartGameHandler(c *gin.Context) {
	sessionID := gs.Sessions.CreateSession()
	c.JSON(http.StatusOK, gin.H{"sessionId": sessionID})
}

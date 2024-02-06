package handlers

import (
	"net/http"

	"github.com/ProlificLabs/captrivia/backend/service"
	"github.com/ProlificLabs/captrivia/backend/utils"
	"github.com/gin-gonic/gin"
)

type AnswerHandler struct {
	router     *gin.Engine
	gameServer *service.GameServer
}

func NewAnswerHandler(r *gin.Engine, gs *service.GameServer) *AnswerHandler {
	return &AnswerHandler{
		router:     r,
		gameServer: gs,
	}
}

func (h *AnswerHandler) RegisterRoutes() {
	h.router.POST("/answer", h.handle)
}

func (h *AnswerHandler) handle(c *gin.Context) {
	var submittedAnswer struct {
		SessionID  string `json:"sessionId"`
		QuestionID string `json:"questionId"`
		Answer     int    `json:"answer"`
	}
	if err := c.ShouldBindJSON(&submittedAnswer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	session, exists := h.gameServer.Sessions.GetSession(submittedAnswer.SessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	correct, err := h.checkAnswer(submittedAnswer.QuestionID, submittedAnswer.Answer)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if correct {
		session.Score += 10 // Increment score for correct answer
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":      correct,
		"currentScore": session.Score, // Return the current score
	})
}

func (h *AnswerHandler) checkAnswer(questionID string, submittedAnswer int) (bool, error) {
	return utils.CheckAnswer(h.gameServer.Questions, questionID, submittedAnswer)
}

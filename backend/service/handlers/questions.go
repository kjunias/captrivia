package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/ProlificLabs/captrivia/backend/model"
	"github.com/ProlificLabs/captrivia/backend/service"
	"github.com/gin-gonic/gin"
)

type QuestionsHandler struct {
	router     *gin.Engine
	gameServer *service.GameServer
}

func NewQuestionsHandler(r *gin.Engine, gs *service.GameServer) *QuestionsHandler {
	return &QuestionsHandler{
		router:     r,
		gameServer: gs,
	}
}

func (h *QuestionsHandler) RegisterRoutes() {
	h.router.GET("/questions", h.handle)
}

func (h *QuestionsHandler) handle(c *gin.Context) {
	shuffledQuestions := shuffleQuestions(h.gameServer.Questions)
	c.JSON(http.StatusOK, shuffledQuestions[:10])
}

func shuffleQuestions(questions []model.Question) []model.Question {
	rand.Seed(time.Now().UnixNano())
	qs := make([]model.Question, len(questions))
	copy(qs, questions)
	rand.Shuffle(len(qs), func(i, j int) {
		qs[i], qs[j] = qs[j], qs[i]
	})
	return qs
}

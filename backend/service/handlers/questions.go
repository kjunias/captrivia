package handlers

import (
	"net/http"

	"github.com/ProlificLabs/captrivia/backend/service"
	"github.com/ProlificLabs/captrivia/backend/utils"
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
	shuffledQuestions := utils.ShuffleQuestions(h.gameServer.Questions)
	c.JSON(http.StatusOK, shuffledQuestions[:10])
}

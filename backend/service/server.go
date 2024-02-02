package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/ProlificLabs/captrivia/backend/model"
	"github.com/ProlificLabs/captrivia/backend/session"
	"github.com/ProlificLabs/captrivia/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type GameRoom struct {
	RoomID               string
	Scores               map[string]int
	AdminID              string
	IsCountingDown       bool
	CurrentQuestionIndex int
	Questions            []model.Question
	mutex                sync.Mutex
}
type GameServer struct {
	Questions []model.Question
	Sessions  *session.SessionStore
	GameRooms map[string]*GameRoom
	Hub       *melody.Melody
}

func NewGameServer(questions []model.Question, store *session.SessionStore) *GameServer {
	return &GameServer{
		Questions: questions,
		Sessions:  store,
		GameRooms: map[string]*GameRoom{},
		Hub:       melody.New(),
	}
}

func (s *GameServer) CreateGameRoom(c *gin.Context) {
	roomID := utils.CreateID(6)
	for _, ok := s.GameRooms[roomID]; ok; {
		roomID = utils.CreateID(6)
	}
	s.GameRooms[roomID] = &GameRoom{
		RoomID: roomID,
		Scores: map[string]int{},
	}

	s.JoinGameRoom(roomID, c, true)
}

func (s *GameServer) JoinGameRoom(roomID string, c *gin.Context, isAdmin bool) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game room not found"})
		return
	}

	playerID := utils.CreateID(6)
	for _, ok := gameRoom.Scores[playerID]; ok; {
		playerID = utils.CreateID(6)
	}
	if isAdmin {
		gameRoom.AdminID = playerID
	}
	gameRoom.Scores[playerID] = 0
	c.JSON(http.StatusOK, getGameRoomResponse(gameRoom, playerID))
	s.publishUpdates(roomID)
}

func (s *GameServer) HandleStartCounter(roomID string, numberOfQuestions int) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		fmt.Println("Error getting game room:", roomID)
		return
	}
	gameRoom.IsCountingDown = true
	gameRoom.CurrentQuestionIndex = -1
	gameRoom.Questions = utils.ShuffleQuestions(s.Questions)[:numberOfQuestions]
	s.publishUpdates(roomID)
}

func (s *GameServer) HandleEndCounter(roomID string) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		fmt.Println("Error getting game room:", roomID)
		return
	}
	gameRoom.IsCountingDown = false
	gameRoom.CurrentQuestionIndex = 0
	s.publishUpdates(roomID)
}

func (s *GameServer) HandleSubmitAnswer(roomID string, playerID string, questionID string, submittedAnswer int) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		fmt.Println("Error getting game room:", roomID)
		return
	}

	gameRoom.mutex.Lock()
	defer gameRoom.mutex.Unlock()
	if questionID != gameRoom.Questions[gameRoom.CurrentQuestionIndex].ID {
		return
	}

	correct, err := utils.CheckAnswer(gameRoom.Questions, questionID, submittedAnswer)
	if err != nil {
		fmt.Println("Error checking answer:", playerID, questionID)
	}

	if correct {
		gameRoom.Scores[playerID] += 10
		gameRoom.CurrentQuestionIndex += 1
		s.publishUpdates(roomID)
	}

	if gameRoom.CurrentQuestionIndex >= len(gameRoom.Questions) {
		delete(s.GameRooms, roomID)
	}
}

func (s *GameServer) StartGameRoom(roomID string, numberOfQuestions int, c *gin.Context) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game room not found"})
		return
	}
	gameRoom.IsCountingDown = true
	gameRoom.Questions = utils.ShuffleQuestions(s.Questions)[:numberOfQuestions]
	c.JSON(http.StatusOK, getGameRoomResponse(gameRoom, gameRoom.AdminID))
	s.publishUpdates(roomID)
}

func (s *GameServer) getGameRoom(roomID string) *GameRoom {
	var gameRoom *GameRoom
	var ok bool
	if gameRoom, ok = s.GameRooms[roomID]; !ok {
		fmt.Println("Game room not found")
		return nil
	}
	return gameRoom
}

type GameRoomResponse struct {
	RoomID               string           `json:"roomID"`
	Scores               map[string]int   `json:"scores"`
	PlayerID             string           `json:"playerID,omitempty"`
	AdminID              string           `json:"adminID"`
	IsCountingDown       bool             `json:"isCountingDown"`
	CurrentQuestionIndex int              `json:"currentQuestionIndex"`
	Questions            []model.Question `json:"questions,omitempty"`
}

func (s *GameServer) publishUpdates(roomID string) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		fmt.Println("Publishing error: Game room not found")
		return
	}
	msg := new(bytes.Buffer)
	if err := json.NewEncoder(msg).Encode(getGameRoomMessage(gameRoom)); err != nil {
		fmt.Println("Failed encoding message...: ", err)
	}
	if err := s.Hub.Broadcast(msg.Bytes()); err != nil {
		fmt.Println("Failed broadcasting message...: ", err)
	}
}

func getGameRoomResponse(gameRoom *GameRoom, playerID string) GameRoomResponse {
	response := getGameRoomMessage(gameRoom)
	response.PlayerID = playerID
	return response
}

func getGameRoomMessage(gameRoom *GameRoom) GameRoomResponse {
	return GameRoomResponse{
		RoomID:               gameRoom.RoomID,
		Scores:               gameRoom.Scores,
		AdminID:              gameRoom.AdminID,
		CurrentQuestionIndex: gameRoom.CurrentQuestionIndex,
		Questions:            gameRoom.Questions,
		IsCountingDown:       gameRoom.IsCountingDown,
	}
}

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ProlificLabs/captrivia/backend/model"
	"github.com/ProlificLabs/captrivia/backend/session"
	"github.com/ProlificLabs/captrivia/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type GameRoomResponse struct {
	RoomID               string           `json:"roomID"`
	Scores               map[string]int   `json:"scores"`
	PlayerID             string           `json:"playerID,omitempty"`
	AdminID              string           `json:"adminID"`
	WinnerID             string           `json:"winnerID,omitempty"`
	CurrentQuestionIndex int              `json:"currentQuestionIndex"`
	Questions            []model.Question `json:"questions,omitempty"`
	State                model.GameState  `json:"state,omitempty"`
}

type GameServer struct {
	Questions []model.Question
	Sessions  *session.SessionStore
	GameRooms map[string]*model.GameRoom
	Hub       *melody.Melody
}

func NewGameServer(questions []model.Question, store *session.SessionStore) *GameServer {
	return &GameServer{
		Questions: questions,
		Sessions:  store,
		GameRooms: map[string]*model.GameRoom{},
		Hub:       melody.New(),
	}
}

func (s *GameServer) CreateGameRoom(c *gin.Context) {
	roomID := utils.CreateID(6)
	for _, ok := s.GameRooms[roomID]; ok; {
		roomID = utils.CreateID(6)
	}
	s.GameRooms[roomID] = &model.GameRoom{
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
		gameRoom.State = model.WAITING
	}
	gameRoom.Scores[playerID] = 0
	c.JSON(http.StatusOK, getGameRoomResponse(gameRoom, playerID))
	s.publishUpdates(gameRoom)
}

func (s *GameServer) HandleStartCounter(roomID string, numberOfQuestions int) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		fmt.Println("Error getting game room:", roomID)
		return
	}
	gameRoom.State = model.COUNTING_DOWN
	gameRoom.CurrentQuestionIndex = 0
	gameRoom.Questions = utils.ShuffleQuestions(s.Questions)[:numberOfQuestions]
	s.publishUpdates(gameRoom)
}

func (s *GameServer) HandleEndCounter(roomID string) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		fmt.Println("Error getting game room:", roomID)
		return
	}
	if gameRoom.State != model.COUNTING_DOWN {
		fmt.Println("Not counting down, discarding end countdown signal")
		return
	}
	gameRoom.State = model.PLAYING
	gameRoom.CurrentQuestionIndex = 0
	s.publishUpdates(gameRoom)
}

func (s *GameServer) HandleSubmitAnswer(roomID string, playerID string, questionID string, submittedAnswer int) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		fmt.Println("Error getting game room:", roomID)
		return
	}

	gameRoom.Mutex.Lock()
	defer gameRoom.Mutex.Unlock()
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
		if gameRoom.CurrentQuestionIndex >= len(gameRoom.Questions) {
			gameRoom.State = model.END
			gameRoom.WinnerID = getWinner(gameRoom)
		}
		s.publishUpdates(gameRoom)
	}
}

func (s *GameServer) StartGameRoom(roomID string, numberOfQuestions int, c *gin.Context) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game room not found"})
		return
	}
	gameRoom.WinnerID = ""
	gameRoom.State = model.COUNTING_DOWN
	gameRoom.Questions = utils.ShuffleQuestions(s.Questions)[:numberOfQuestions]
	c.JSON(http.StatusOK, getGameRoomResponse(gameRoom, gameRoom.AdminID))
	s.publishUpdates(gameRoom)
}

func (s *GameServer) getGameRoom(roomID string) *model.GameRoom {
	var gameRoom *model.GameRoom
	var ok bool
	if gameRoom, ok = s.GameRooms[roomID]; !ok {
		fmt.Println("Game room not found")
		return nil
	}
	return gameRoom
}

func (s *GameServer) publishUpdates(gameRoom *model.GameRoom) {
	msg := new(bytes.Buffer)
	update := getGameRoomMessage(gameRoom)
	if err := json.NewEncoder(msg).Encode(update); err != nil {
		fmt.Println("Failed encoding message...: ", err)
	}
	if err := s.Hub.Broadcast(msg.Bytes()); err != nil {
		fmt.Println("Failed broadcasting message...: ", err)
	}
}

func getWinner(gameRoom *model.GameRoom) string {
	winner := ""
	maxScore := -1
	for p, s := range gameRoom.Scores {
		if maxScore < s {
			winner = p
			maxScore = s
		}
	}
	return winner
}

func getGameRoomResponse(gameRoom *model.GameRoom, playerID string) GameRoomResponse {
	response := getGameRoomMessage(gameRoom)
	response.PlayerID = playerID
	return response
}

func getGameRoomMessage(gameRoom *model.GameRoom) GameRoomResponse {
	return GameRoomResponse{
		RoomID:               gameRoom.RoomID,
		Scores:               gameRoom.Scores,
		AdminID:              gameRoom.AdminID,
		WinnerID:             gameRoom.WinnerID,
		CurrentQuestionIndex: gameRoom.CurrentQuestionIndex,
		Questions:            gameRoom.Questions,
		State:                gameRoom.State,
	}
}

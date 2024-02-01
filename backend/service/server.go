package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ProlificLabs/captrivia/backend/model"
	"github.com/ProlificLabs/captrivia/backend/session"
	"github.com/ProlificLabs/captrivia/backend/utils"
	"github.com/gin-gonic/gin"
)

type GameRoom struct {
	RoomID          string
	Channels        map[string]chan *gin.Context
	Scores          map[string]int
	AdminID         string
	IsCountingDown  bool
	CurrentQuestion int
	Questions       []model.Question
}
type GameServer struct {
	Questions []model.Question
	Sessions  *session.SessionStore
	GameRooms map[string]*GameRoom
}

func NewGameServer(questions []model.Question, store *session.SessionStore) *GameServer {
	return &GameServer{
		Questions: questions,
		Sessions:  store,
		GameRooms: map[string]*GameRoom{},
	}
}

func (s *GameServer) CreateGameRoom(c *gin.Context) {
	fmt.Printf("\n\n ====CreateGameRoom %#v\n", "====")
	roomID := utils.CreateID(6)
	for _, ok := s.GameRooms[roomID]; ok; {
		roomID = utils.CreateID(6)
	}
	s.GameRooms[roomID] = &GameRoom{
		RoomID:   roomID,
		Channels: map[string]chan *gin.Context{},
		Scores:   map[string]int{},
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
	fmt.Println("==> GameServe.JoinGameRoom playerID 0: ", playerID)
	for _, ok := gameRoom.Channels[playerID]; ok; {
		playerID = utils.CreateID(6)
	}
	fmt.Println("==> GameServe.JoinGameRoom playerID 1: ", playerID)
	// fmt.Printf("\n\n ==> playerID: %#v\n", playerID)
	fmt.Println("==> GameServe.JoinGameRoom: ")
	if isAdmin {
		gameRoom.AdminID = playerID
	}
	gameRoom.Scores[playerID] = 0
	ch := make(chan *gin.Context, 1)
	ch <- c
	gameRoom.Channels[playerID] = ch
	s.publishUpdates(roomID)
}

func (s *GameServer) SubscribeToGameRoom(roomID string, playerID string, c *gin.Context) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game room not found"})
		// c.Request.Context().Done()
		return
	}
	playerChannel := gameRoom.getPlayerChannel(playerID)
	playerChannel <- c
	ok := false
	for !ok {
		select {
		case playerChannel <- c:
			ok = true
		default:
			fmt.Println("Channel full... flushing and retrying")
			(<-playerChannel).Request.Context().Done()
		}
	}
	fmt.Println("==> GameServer.SubscribeToGameRoom: ", roomID, playerID)
}

func (s *GameServer) StartGameRoom(roomID string, numberOfQuestions int, c *gin.Context) {
	fmt.Println("=====> GameServer.StartGameRoom start 1")
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game room not found"})
		return
	}
	gameRoom.IsCountingDown = true
	gameRoom.Questions = utils.ShuffleQuestions(s.Questions)[:numberOfQuestions]
	fmt.Println("===> GameServe.StartGameRoom: ", gameRoom)
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

func (r *GameRoom) getPlayerChannel(playerID string) chan *gin.Context {
	var channel chan *gin.Context
	var ok bool
	if channel, ok = r.Channels[playerID]; !ok {
		log.Fatal("Player channle not found")
	}
	return channel
}

type GameRoomResponse struct {
	RoomID          string           `json:"roomID"`
	Scores          map[string]int   `json:"scores"`
	PlayerID        string           `json:"playerID"`
	AdminID         string           `json:"adminID"`
	IsCountingDown  bool             `json:"isCountingDown"`
	CurrentQuestion int              `json:"currentQuestion"`
	Questions       []model.Question `json:"questions"`
}

func (s *GameServer) publishUpdates(roomID string) {
	gameRoom := s.getGameRoom(roomID)
	if gameRoom == nil {
		fmt.Println("Publishing error: Game room not found")
		return
	}
	for id, c := range gameRoom.Channels {
		fmt.Println("==> publishUpdates: ", id)
		empty := false
		for !empty {
			select {
			case ctx := <-c:
				fmt.Println("===> publishing...: ", id)
				ctx.JSON(http.StatusOK, getGameRoomResponse(gameRoom, id))
				if ctx != nil && ctx.Writer != nil {
					ctx.Writer.Flush()
					ctx.Request.Context().Done()
				}
			default:
				empty = true
			}
		}
	}
}

func getGameRoomResponse(gameRoom *GameRoom, playerID string) GameRoomResponse {
	return GameRoomResponse{
		RoomID:          gameRoom.RoomID,
		PlayerID:        playerID,
		Scores:          gameRoom.Scores,
		AdminID:         gameRoom.AdminID,
		CurrentQuestion: gameRoom.CurrentQuestion,
		Questions:       gameRoom.Questions,
		IsCountingDown:  gameRoom.IsCountingDown,
	}
}

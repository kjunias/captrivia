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
	Channels map[string]chan *gin.Context
	Scores   map[string]int
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

func (s *GameServer) CreateGameRoom(c *gin.Context) (string, string) {
	fmt.Printf("\n\n ====CreateGameRoom %#v\n", "====")
	roomID := utils.CreateID(6)
	for _, ok := s.GameRooms[roomID]; ok; {
		roomID = utils.CreateID(6)
	}
	fmt.Printf("\n\n ==> roomID: %#v\n", roomID)

	s.GameRooms[roomID] = &GameRoom{
		Channels: map[string]chan *gin.Context{},
		Scores:   map[string]int{},
	}
	return s.JoinGameRoom(roomID, c)
}

func (s *GameServer) JoinGameRoom(roomID string, c *gin.Context) (string, string) {
	gameRoom := s.getGameRoom(roomID)
	playerID := utils.CreateID(6)
	fmt.Println("==> GameServe.JoinGameRoom playerID 0: ", playerID)
	for _, ok := gameRoom.Channels[playerID]; ok; {
		playerID = utils.CreateID(6)
	}
	fmt.Println("==> GameServe.JoinGameRoom playerID 1: ", playerID)
	// fmt.Printf("\n\n ==> playerID: %#v\n", playerID)
	fmt.Println("==> GameServe.JoinGameRoom: ")
	gameRoom.Scores[playerID] = 0
	ch := make(chan *gin.Context, 1)
	ch <- c
	gameRoom.Channels[playerID] = ch
	s.publishUpdates(roomID)
	return roomID, playerID
}

func (s *GameServer) SubscribeToGameRoom(roomID string, playerID string, c *gin.Context) {
	gameRoom := s.getGameRoom(roomID)
	playerChannel := gameRoom.getPlayerChannel(playerID)
	playerChannel <- c
	// ok := false
	// for !ok {
	// select {
	// case playerChannel <- c:
	// ok = true
	// default:
	// fmt.Println("Channel full... flushing and retrying")
	// (<-playerChannel)
	// }
	// }
	fmt.Println("==> GameServer.SubscribeToGameRoom: ", roomID, playerID)
}

func (s *GameServer) getGameRoom(roomID string) *GameRoom {
	var gameRoom *GameRoom
	var ok bool
	if gameRoom, ok = s.GameRooms[roomID]; !ok {
		log.Fatal("Game room not found")
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

type ResponseLou struct {
	RoomID   string         `json:"roomID"`
	Scores   map[string]int `json:"scores"`
	PlayerID string         `json:"playerID"`
}

func (gr *GameServer) publishUpdates(roomID string) {
	gameRoom := gr.getGameRoom(roomID)
	for id, c := range gameRoom.Channels {
		fmt.Println("==> publishUpdates: ", id)
		ctx := <-c

		dat := ResponseLou{
			RoomID:   roomID,
			PlayerID: id,
			Scores:   gameRoom.Scores,
		}

		ctx.JSON(http.StatusOK, dat)
		// close(c)
	}
}

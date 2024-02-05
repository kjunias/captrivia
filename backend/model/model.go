package model

import "sync"

type GameState string

const (
	WAITING       GameState = "WAITING"
	COUNTING_DOWN           = "COUNTING_DOWN"
	PLAYING                 = "PLAYING"
	END                     = "END"
)

type GameRoom struct {
	RoomID               string
	Scores               map[string]int
	AdminID              string
	WinnerID             string
	CurrentQuestionIndex int
	State                GameState
	Questions            []Question
	Mutex                sync.Mutex
}

type Question struct {
	ID           string   `json:"id"`
	QuestionText string   `json:"questionText"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
}

type PlayerSession struct {
	Score int
}

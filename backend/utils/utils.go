package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/ProlificLabs/captrivia/backend/model"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func CreateID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ShuffleQuestions(questions []model.Question) []model.Question {
	qs := make([]model.Question, len(questions))
	copy(qs, questions)
	rand.Shuffle(len(qs), func(i, j int) {
		qs[i], qs[j] = qs[j], qs[i]
	})
	return qs
}

func GetIntFromMessage(prop string, msg []byte) (int, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(msg, &m)
	if err != nil {
		fmt.Println("Error unmarshalling message:", msg)
		return 0, err
	}
	p, ok := m[prop].(string)
	if !ok {
		return 0, fmt.Errorf("Error getting property")
	}
	r, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println("Error converting string to int:", r)
		return 0, err
	}
	return r, nil
}

func GetStringFromMessage(prop string, msg []byte) (string, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(msg, &m)
	if err != nil {
		fmt.Println("Error unmarshalling message:", msg)
		return "", err
	}
	p, ok := m[prop].(string)
	if !ok {
		return "", fmt.Errorf("Error getting property")
	}
	return p, nil
}

func CheckAnswer(questions []model.Question, questionID string, submittedAnswer int) (bool, error) {
	for _, question := range questions {
		if question.ID == questionID {
			return question.CorrectIndex == submittedAnswer, nil
		}
	}
	return false, errors.New("question not found")
}

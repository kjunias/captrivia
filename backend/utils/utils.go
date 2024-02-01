package utils

import (
	"math/rand"

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

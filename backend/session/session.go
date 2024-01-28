package session

import (
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/ProlificLabs/captrivia/backend/model"
)

type SessionStore struct {
	sync.Mutex
	Sessions map[string]*model.PlayerSession
}

func (store *SessionStore) CreateSession() string {
	store.Lock()
	defer store.Unlock()

	uniqueSessionID := generateSessionID()
	store.Sessions[uniqueSessionID] = &model.PlayerSession{Score: 0}

	return uniqueSessionID
}

func generateSessionID() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func (store *SessionStore) GetSession(sessionID string) (*model.PlayerSession, bool) {
	store.Lock()
	defer store.Unlock()

	session, exists := store.Sessions[sessionID]
	return session, exists
}

package session

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SessionID string

type sessionField struct {
	expiry    int64
	authState string
}
type Session map[SessionID]*sessionField

func (s Session) ReadSession(sessionID SessionID) (SessionID, bool) {
	// No session or session expired
	if t, ok := s[sessionID]; !ok {
		return s.createSession(), false
	} else if t.expiry < time.Now().Unix() {
		return s.createSession(), false
	}
	return sessionID, true
}

func (s Session) createSession() SessionID {
	sessionID := SessionID(uuid.New().String())
	s[sessionID] = &sessionField{expiry: time.Now().AddDate(0, 0, 1).Unix()}
	return SessionID(sessionID)
}

func (s Session) AddStateToSession(id SessionID, state string) {
	fmt.Printf("sessionID: %s \n", id)
	field, ok := s[id]
	if !ok {
		fmt.Println("session not found")
		return
	}
	field.authState = state
	s[id] = field
}

func (s Session) CompareState(id SessionID, state string) bool {
	field := s[id]
	fmt.Printf("sessionID: %s \n", id)
	fmt.Printf("incomingState: %s \n", state)
	fmt.Printf("storageState: %s \n", field.authState)
	return strings.EqualFold(field.authState, state)
}

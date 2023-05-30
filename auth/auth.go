package auth

import "go_project.com/m/session"

type TokenData struct {
	AccessToken  string
	RefreshToken string
}

type AuthStorage map[session.SessionID]*TokenData

func (a AuthStorage) IsAuthenticated(sessionID session.SessionID) bool {
	_, ok := a[sessionID]
	return ok
}

func (a AuthStorage) AddToken(sessionID session.SessionID, token TokenData) {
	a[sessionID] = &token
}

func (a AuthStorage) GetToken(sessionID session.SessionID) string {
	token, ok := a[sessionID]
	if !ok {
		return ""
	}
	return token.AccessToken
}

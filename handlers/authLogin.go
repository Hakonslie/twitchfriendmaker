package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/hakonslie/twitchfriendmaker/config"
	"github.com/hakonslie/twitchfriendmaker/logger"
	"github.com/hakonslie/twitchfriendmaker/session"
)

func generateRandomState(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)[:length]
}

func AuthLogin(sessions *session.Session, cfg *config.Config, log logger.Logger) gin.HandlerFunc {
	baseurl := &url.URL{
		Scheme: "https",
		Host:   "id.twitch.tv",
	}
	path := "/oauth2/authorize"
	state := generateRandomState(32)

	queryParams := url.Values{
		"response_type": {"code"},
		"client_id":     {cfg.ClientID},
		"redirect_uri":  {cfg.RedirectURL},
		"scope":         {"user:read:follows"},
		"state":         {state},
	}
	return func(c *gin.Context) {
		sessionID := c.MustGet("sessionID").(session.SessionID)
		sessions.AddStateToSession(sessionID, state)
		u := baseurl.ResolveReference(&url.URL{Path: path, RawQuery: queryParams.Encode()})
		c.Redirect(http.StatusSeeOther, u.String())

	}
}

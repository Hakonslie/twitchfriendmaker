package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hakonslie/twitchfriendmaker/auth"
	"github.com/hakonslie/twitchfriendmaker/config"
	"github.com/hakonslie/twitchfriendmaker/session"
)

func UseSessionMiddleware(sessions *session.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("sessionID")
		if err != nil {
			// No cookie
			fmt.Printf("no cookie \n")
			newSessionID, _ := sessions.ReadSession("")
			c.SetCookie("sessionID", string(newSessionID), 3600, "/", "localhost", false, true)
		} else {
			// check if cookie is still valid, if not set new
			newSessionID, valid := sessions.ReadSession(session.SessionID(sessionID))
			if !valid {
				c.SetCookie("sessionID", string(newSessionID), 3600, "/", "localhost", false, true)
			}
		}
		c.Set("sessionID", session.SessionID(sessionID))
		c.Next()
	}
}

func UseAuthMiddleware(cfg *config.Config, auth *auth.AuthStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path != "/auth" && c.Request.URL.Path != "/redirect" && c.Request.URL.Path != "/noauth" {
			sessionID := c.MustGet("sessionID").(session.SessionID)
			if !auth.IsAuthenticated(sessionID) {
				c.Redirect(http.StatusSeeOther, cfg.NoAuthURL)
				return
			}
		}
		c.Next()
	}
}

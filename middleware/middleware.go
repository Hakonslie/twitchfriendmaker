package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hakonslie/twitchfriendmaker/auth"
	"github.com/hakonslie/twitchfriendmaker/config"
	"github.com/hakonslie/twitchfriendmaker/session"
	lru "github.com/hashicorp/golang-lru/v2"
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

func UseRateLimiterMiddleware() gin.HandlerFunc {

	type clientRateLimit struct {
		count int
		start int64
	}

	Arc, err := lru.NewARC[string, clientRateLimit](50) // 100 IP-s stored
	if err != nil {
		fmt.Println(err)
	}

	return func(c *gin.Context) {
		key := c.ClientIP()
		if entries, ok := Arc.Get(key); !ok {
			entry := clientRateLimit{
				count: 1,
				start: time.Now().Unix(),
			}
			Arc.Add(key, entry)
			c.Next()
			return
		} else if entries.count > 5 && time.Now().Unix()-entries.start <= 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		} else if time.Now().Unix()-entries.start > 5 {
			entry := clientRateLimit{
				count: 1,
				start: time.Now().Unix(),
			}
			Arc.Add(key, entry)
			c.Next()
			return
		} else {
			// Increment visits
			entry := clientRateLimit{
				count: entries.count + 1,
				start: entries.start,
			}
			fmt.Println(entry)
			Arc.Add(key, entry)
			c.Next()
			return
		}
	}
}

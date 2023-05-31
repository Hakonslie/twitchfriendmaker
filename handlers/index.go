package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hakonslie/twitchfriendmaker/follows"
	"github.com/hakonslie/twitchfriendmaker/logger"
	"github.com/hakonslie/twitchfriendmaker/session"
)

func Index(follows follows.FollowStorage, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.MustGet("sessionID").(session.SessionID)
		follows := follows.GetFollows(sessionID)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":   "Find meaningful relationships",
			"follows": follows,
		})
	}
}

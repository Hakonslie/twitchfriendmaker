package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go_project.com/m/follows"
	"go_project.com/m/logger"
	"go_project.com/m/session"
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

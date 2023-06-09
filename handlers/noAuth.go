package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "noauth.tmpl", gin.H{
			"title": "Welcome",
		})
	}
}

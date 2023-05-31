package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(UseRateLimiterMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Test endpoint"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	for i := 0; i <= 5; i++ {
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("Expected HTTP status OK, got: %d", rec.Code)
		}
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Expected HTTP status 429 Too Many Requests, got: %d", rec.Code)
	}

	// wait for rateLimiter to reset
	time.Sleep(time.Second * 6)

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected HTTP status OK, got: %d", rec.Code)
	}

	for i := 0; i < 5; i++ {
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("Expected HTTP status OK, got: %d", rec.Code)
		}
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hakonslie/twitchfriendmaker/auth"
	"github.com/hakonslie/twitchfriendmaker/config"
	"github.com/hakonslie/twitchfriendmaker/follows"
	"github.com/hakonslie/twitchfriendmaker/handlers"
	"github.com/hakonslie/twitchfriendmaker/logger"
	"github.com/hakonslie/twitchfriendmaker/middleware"
	"github.com/hakonslie/twitchfriendmaker/session"
)

type keyValue struct {
	key   string
	value string
}

func runServer(ctx context.Context, log logger.Logger, wg *sync.WaitGroup, config *config.Config, router *gin.Engine) {
	defer wg.Done()
	server := &http.Server{
		Addr:    config.Port,
		Handler: router,
	}
	wg.Add(1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("Server error")
			fmt.Println(err)
		}
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			server.Shutdown(ctx)
			return
		}
	}
}
func setupRouter(log logger.Logger, cfg *config.Config) *gin.Engine {
	sessions := session.Session{}
	auth := auth.AuthStorage{}
	follows := follows.FollowStorage{}
	router := gin.Default()
	router.StaticFile("/styles.css", "../templates/styles.css")
	router.Use(middleware.UseSessionMiddleware(&sessions))
	router.Use(middleware.UseAuthMiddleware(cfg, &auth))

	router.LoadHTMLGlob("../templates/*")
	router.GET("/", handlers.Index(follows, log))
	router.GET("/redirect", handlers.Redirect(&sessions, &auth, cfg, log))
	router.GET("/auth", handlers.AuthLogin(&sessions, cfg, log))
	router.GET("/noauth", handlers.NoAuth())
	router.GET("/getFollows", handlers.GetFollows(follows, auth, cfg))

	return router
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.SetupLogger(ctx, "../logs/")
	router := setupRouter(log, &cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go runServer(ctx, log, &wg, &cfg, router)
	wg.Add(1)
	go func() {
		for {
			select {
			case <-quit:
				log.Info("Shutting down")
				cancel()
				wg.Done()
			}
		}
	}()
	wg.Wait()
	log.Info("Done")
}

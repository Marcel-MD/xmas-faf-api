package handlers

import (
	"os"
	"sync"

	"github.com/Marcel-MD/xmas-faf-api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var once sync.Once

func InitRouter() {
	once.Do(func() {
		log.Info().Msg("Initializing router")

		e := gin.Default()

		e.Use(middleware.CORS(), middleware.RateLimiter())

		r := e.Group("/api")

		routeUserHandler(r)
		routeTrainingHandler(r)
		routePostHandler(r)
		routeFileHandler(r)
		routeCommentHandler(r)

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		e.Run(":" + port)
	})
}

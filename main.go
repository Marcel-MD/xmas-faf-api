package main

import (
	"github.com/Marcel-MD/xmas-faf-api/handlers"
	"github.com/Marcel-MD/xmas-faf-api/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Err(err).Msg("Error loading .env file")
	}

	logger.Config()
	handlers.InitRouter()
}

package main

import (
	"log/slog"
	_ "net/http/pprof"
	"os"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"github.com/EdwinRincon/browersfc-api/server"
	"github.com/joho/godotenv"
)

// @BasePath	/api
// @title BrowersFC API
// @version 1.0
// @description API para la gestión de la liga de fútbol BrowersFC
func main() {
	// Load environment variables for local development
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(".env"); err != nil {
			slog.Warn("Error loading .env file, relying on environment variables", "error", err)
		}
	}

	// Initialize configurations
	config.InitLogConfig()

	// Setup centralized logger
	logger.Setup(logger.LogConfig{
		Level:   config.AppLogConfig.Level,
		Format:  logger.LogFormat(config.AppLogConfig.Format),
		Output:  os.Stdout,
		IsDebug: os.Getenv("GIN_MODE") != "release",
	})

	// Initialize OAuth
	if err := config.InitOAuth(); err != nil {
		slog.Error("Failed to initialize OAuth", slog.String("error", err.Error()))
		os.Exit(1)
	}

	server := server.NewServer()
	server.Start()
}

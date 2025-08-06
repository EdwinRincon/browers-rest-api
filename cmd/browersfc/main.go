package main

import (
	"log/slog"
	_ "net/http/pprof"
	"os"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/server"
	"github.com/joho/godotenv"
)

// @BasePath	/api
// @title BrowersFC API
// @version 1.0
// @description API para la gestión de la liga de fútbol BrowersFC
func main() {
	// Setup structured logger
	logLevel := slog.LevelInfo
	if os.Getenv("GIN_MODE") != "release" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Load environment variables for local development
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(".env"); err != nil {
			// This is a warning because in a container, env vars are injected directly.
			slog.Warn("Error loading .env file, relying on environment variables", "error", err)
		}
	}
	// Initialize OAuth
	if err := config.InitOAuth(); err != nil {
		slog.Error("Failed to initialize OAuth", slog.String("error", err.Error()))
		os.Exit(1)
	}

	server := server.NewServer()
	server.Start()
}

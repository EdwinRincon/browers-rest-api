package main

import (
	"flag"
	"log/slog"
	_ "net/http/pprof"
	"os"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"github.com/EdwinRincon/browersfc-api/pkg/seed"
	"github.com/EdwinRincon/browersfc-api/server"
	"github.com/joho/godotenv"
)

// @BasePath	/api
// @title BrowersFC API
// @version 1.0
// @description API para la gestión de la liga de fútbol BrowersFC
func main() {
	// Only seed the database if the flag is set
	// dev only
	seedDB := flag.Bool("seed", false, "Seed the database with initial data")
	flag.Parse()

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

	// Seed database if requested
	if *seedDB {
		slog.Info("Seeding database...")
		if err := seed.SeedDatabase(); err != nil {
			slog.Error("Failed to seed database", slog.String("error", err.Error()))
			os.Exit(1)
		}
		slog.Info("Database seeded successfully")
		return
	}

	if err := config.InitOAuth(); err != nil {
		slog.Error("Failed to initialize OAuth", slog.String("error", err.Error()))
		os.Exit(1)
	}

	server := server.NewServer()
	server.Start()
}

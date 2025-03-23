package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables first
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize OAuth
	if err := config.InitOAuth(); err != nil {
		log.Fatalf("Failed to initialize OAuth: %v", err)
	}

	server := server.NewServer()
	server.Start()
}

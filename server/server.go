package server

import (
	"log"
	"net/http"
	"time"

	"github.com/EdwinRincon/browersfc-api/api"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/repository"
	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/pkg/orm"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Server struct {
	Router    *gin.Engine
	Port      string
	JWTSecret string
}

func NewServer() *Server {
	// Carga variables de entorno desde .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error al cargar variables de entorno desde .env ->" + err.Error())
	}
	// Configura el enrutador Gin
	r := gin.Default()

	// Configura la base de datos MySQL
	db, err := orm.NewConnectionDB()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	api.InitializeRoutes(r, userHandler)

	// Modificación para leer el puerto desde un archivo
	port := config.GetPort()

	// Modificación para leer el JWTSecret desde un archivo
	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		log.Fatalf("Failed to read JWT secret from file: %v", err)
	}
	if jwtSecret == "" {
		log.Fatal("JWT secret is not defined")
	}

	return &Server{
		Router:    r,
		Port:      port,
		JWTSecret: jwtSecret,
	}
}

// Start method is modified to use http.Server for more detailed configuration.
func (s *Server) Start() {
	server := &http.Server{
		Addr:    ":" + s.Port,
		Handler: s.Router,

		// Set timeouts to avoid Slowloris attacks
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", s.Port, err)
	}
}

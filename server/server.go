package server

import (
	"log"
	"net/http"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/repository"
	router "github.com/EdwinRincon/browersfc-api/api/router"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
	"github.com/EdwinRincon/browersfc-api/pkg/orm"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Server struct {
	Router    *gin.Engine
	Port      string
	JWTSecret []byte
}

// NewServer configura y devuelve una instancia del servidor.
func NewServer() *Server {
	// Carga variables de entorno desde .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	// Configura el enrutador Gin
	r := gin.Default()

	// Leer el puerto desde un archivo de configuración
	port := config.GetPort()

	// Leer el JWTSecret desde un archivo de configuración
	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		log.Fatalf("Failed to read JWT secret from file: %v", err)
	}
	if jwtSecret == nil {
		log.Fatal("JWT secret is not defined")
	}

	return &Server{
		Router:    r,
		Port:      port,
		JWTSecret: jwtSecret,
	}
}

// SetupRouter configura las rutas del servidor.
func (s *Server) SetupRouter() {
	db, err := orm.NewConnectionDB()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Inicializar los repositorios
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Inicializar los servicios
	jwtService := jwt.NewJWTService(string(s.JWTSecret))
	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)

	// Inicializar los manejadores
	userHandler := handler.NewUserHandler(authService, userService)
	roleHandler := handler.NewRoleHandler(roleService)

	// Configurar las rutas
	router.InitializeUserRoutes(s.Router, userHandler)
	router.InitializeRoleRoutes(s.Router, roleHandler)
}

// Start inicia el servidor con la configuración de HTTP detallada.
func (s *Server) Start() {
	// Configurar las rutas antes de iniciar el servidor
	s.SetupRouter()

	// Configurar el servidor HTTP con tiempos de espera para seguridad
	server := &http.Server{
		Addr:              ":" + s.Port,
		Handler:           s.Router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Iniciar el servidor y manejar errores de inicio
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", s.Port, err)
	}
}

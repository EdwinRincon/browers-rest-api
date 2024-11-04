package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/EdwinRincon/browersfc-api/api/repository"
	router "github.com/EdwinRincon/browersfc-api/api/router"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/config"
	docs "github.com/EdwinRincon/browersfc-api/docs"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
	"github.com/EdwinRincon/browersfc-api/pkg/orm"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath	/api
// @title BrowersFC API
// @version 1.0
// @description API para la gestión de la liga de fútbol BrowersFC

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

	r.Use(middleware.RateLimit())

	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "BrowersFC API"
	docs.SwaggerInfo.Description = "API para la gestión de la liga de fútbol BrowersFC"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"https"}
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
	db, err := orm.GetDBInstance()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
		return
	}

	// Inicializar los repositorios
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	seasonRepo := repository.NewSeasonRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	classificationRepo := repository.NewClassificationRepository(db)
	lineupRepo := repository.NewLineupRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	teamStatsRepo := repository.NewTeamStatsRepository(db)

	// Inicializar los servicios
	jwtService := jwt.NewJWTService(string(s.JWTSecret))
	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	teamService := service.NewTeamService(teamRepo)
	playerService := service.NewPlayerService(playerRepo)
	seasonService := service.NewSeasonService(seasonRepo)
	articleService := service.NewArticleService(articleRepo)
	classificationService := service.NewClassificationService(classificationRepo)
	lineupService := service.NewLineupService(lineupRepo)
	matchService := service.NewMatchService(matchRepo)
	teamStatsService := service.NewTeamStatsService(teamStatsRepo)

	// Inicializar los manejadores
	userHandler := handler.NewUserHandler(authService, userService)
	roleHandler := handler.NewRoleHandler(roleService)
	teamHandler := handler.NewTeamHandler(teamService)
	playerHandler := handler.NewPlayerHandler(playerService)
	seasonHandler := handler.NewSeasonHandler(seasonService)
	articleHandler := handler.NewArticleHandler(articleService)
	classificationHandler := handler.NewClassificationHandler(classificationService)
	lineupHandler := handler.NewLineupHandler(lineupService)
	matchHandler := handler.NewMatchHandler(matchService)
	teamStatsHandler := handler.NewTeamStatsHandler(teamStatsService)

	// Configurar las rutas
	router.InitializeUserRoutes(s.Router, userHandler)
	router.InitializeRoleRoutes(s.Router, roleHandler)
	router.InitializeTeamRoutes(s.Router, teamHandler)
	router.InitializePlayerRoutes(s.Router, playerHandler)
	router.InitializeSeasonRoutes(s.Router, seasonHandler)
	router.InitializeArticleRoutes(s.Router, articleHandler)
	router.InitializeClassificationRoutes(s.Router, classificationHandler)
	router.InitializeLineupRoutes(s.Router, lineupHandler)
	router.InitializeMatchRoutes(s.Router, matchHandler)
	router.InitializeTeamStatsRoutes(s.Router, teamStatsHandler)

	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", s.Port, err)
		}
	}()

	// Channel that listens for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // SIGINT for Ctrl+C, SIGTERM for termination
	<-quit                                             // Wait until we receive a termination signal

	log.Println("Shutting down server...")

	// Context with timeout for graceful shutdown (e.g., 10 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

}

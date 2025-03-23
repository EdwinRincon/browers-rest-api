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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
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
	// Configura el enrutador Gin
	r := gin.Default()
	r.Use(middleware.RateLimit())

	// Configurar Swagger
	setupSwagger()

	// Leer el puerto desde un archivo de configuración
	port := config.GetPort()

	// Leer el JWTSecret desde un archivo de configuración
	jwtSecret := getJWTSecret()

	// Iniciar un servidor HTTP adicional para el profiling solo en desarrollo
	setupPprofServer()

	return &Server{
		Router:    r,
		Port:      port,
		JWTSecret: jwtSecret,
	}
}

// SetupRouter configura las rutas del servidor.
func (s *Server) SetupRouter() {
	// Obtener conexión a la base de datos
	db := getDBConnection()

	// Inicializar componentes
	repositories := initializeRepositories(db)
	services := initializeServices(repositories, s.JWTSecret)
	handlers := initializeHandlers(services)

	// Configurar las rutas
	initializeRoutes(s.Router, handlers)

	// Configurar Swagger
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Start inicia el servidor con la configuración de HTTP detallada.
func (s *Server) Start() {
	// Configurar las rutas antes de iniciar el servidor
	s.SetupRouter()

	// Configurar el servidor HTTP con tiempos de espera para seguridad
	server := createHTTPServer(s.Port, s.Router)

	// Iniciar el servidor en una goroutine
	go startServer(server, s.Port)

	// Esperar señal de apagado y realizar shutdown graceful
	gracefulShutdown(server)
}

// Funciones auxiliares

func setupSwagger() {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "BrowersFC API"
	docs.SwaggerInfo.Description = "API para la gestión de la liga de fútbol BrowersFC"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"https"}
}

func getJWTSecret() []byte {
	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		log.Fatalf("Failed to read JWT secret from file: %v", err)
	}
	if jwtSecret == nil {
		log.Fatal("JWT secret is not defined")
	}
	return jwtSecret
}

func setupPprofServer() {
	if os.Getenv("GIN_MODE") != "release" {
		go func() {
			log.Println("Starting pprof server on :6060")
			if err := http.ListenAndServe(":6060", nil); err != nil {
				log.Fatalf("Could not start pprof server: %v", err)
			}
		}()
	}
}

func getDBConnection() *gorm.DB {
	db, err := orm.GetDBInstance()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	return db
}

func createHTTPServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func startServer(server *http.Server, port string) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}
}

func gracefulShutdown(server *http.Server) {
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

// Funciones para inicializar componentes

func initializeRepositories(db *gorm.DB) map[string]interface{} {
	return map[string]interface{}{
		"user":      repository.NewUserRepository(db),
		"role":      repository.NewRoleRepository(db),
		"team":      repository.NewTeamRepository(db),
		"player":    repository.NewPlayerRepository(db),
		"season":    repository.NewSeasonRepository(db),
		"article":   repository.NewArticleRepository(db),
		"lineup":    repository.NewLineupRepository(db),
		"match":     repository.NewMatchRepository(db),
		"teamStats": repository.NewTeamStatsRepository(db),
	}
}

func initializeServices(repos map[string]interface{}, jwtSecret []byte) map[string]interface{} {
	jwtService := jwt.NewJWTService(string(jwtSecret))

	return map[string]interface{}{
		"jwt":       jwtService,
		"auth":      service.NewAuthService(repos["user"].(repository.UserRepository), jwtService),
		"user":      service.NewUserService(repos["user"].(repository.UserRepository)),
		"role":      service.NewRoleService(repos["role"].(repository.RoleRepository)),
		"team":      service.NewTeamService(repos["team"].(repository.TeamRepository)),
		"player":    service.NewPlayerService(repos["player"].(repository.PlayerRepository)),
		"season":    service.NewSeasonService(repos["season"].(repository.SeasonRepository)),
		"article":   service.NewArticleService(repos["article"].(repository.ArticleRepository)),
		"lineup":    service.NewLineupService(repos["lineup"].(repository.LineupRepository)),
		"match":     service.NewMatchService(repos["match"].(repository.MatchRepository)),
		"teamStats": service.NewTeamStatsService(repos["teamStats"].(repository.TeamStatsRepository)),
	}
}

func initializeHandlers(services map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"user":      handler.NewUserHandler(services["auth"].(service.AuthService), services["user"].(service.UserService)),
		"role":      handler.NewRoleHandler(services["role"].(service.RoleService)),
		"team":      handler.NewTeamHandler(services["team"].(service.TeamService)),
		"player":    handler.NewPlayerHandler(services["player"].(service.PlayerService)),
		"season":    handler.NewSeasonHandler(services["season"].(service.SeasonService)),
		"article":   handler.NewArticleHandler(services["article"].(service.ArticleService)),
		"lineup":    handler.NewLineupHandler(services["lineup"].(service.LineupService)),
		"match":     handler.NewMatchHandler(services["match"].(service.MatchService)),
		"teamStats": handler.NewTeamStatsHandler(services["teamStats"].(service.TeamStatsService)),
	}
}

func initializeRoutes(r *gin.Engine, handlers map[string]interface{}) {
	router.InitializeUserRoutes(r, handlers["user"].(*handler.UserHandler))
	router.InitializeRoleRoutes(r, handlers["role"].(*handler.RoleHandler))
	router.InitializeTeamRoutes(r, handlers["team"].(*handler.TeamHandler))
	router.InitializePlayerRoutes(r, handlers["player"].(*handler.PlayerHandler))
	router.InitializeSeasonRoutes(r, handlers["season"].(*handler.SeasonHandler))
	router.InitializeArticleRoutes(r, handlers["article"].(*handler.ArticleHandler))
	router.InitializeLineupRoutes(r, handlers["lineup"].(*handler.LineupHandler))
	router.InitializeMatchRoutes(r, handlers["match"].(*handler.MatchHandler))
	router.InitializeTeamStatsRoutes(r, handlers["teamStats"].(*handler.TeamStatsHandler))
}

package server

import (
	"context"
	"log/slog"
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

type Server struct {
	Router    *gin.Engine
	Port      string
	JWTSecret []byte
}

// Dependency containers for type-safe injection
type Repositories struct {
	User      repository.UserRepository
	Role      repository.RoleRepository
	Team      repository.TeamRepository
	Player    repository.PlayerRepository
	Season    repository.SeasonRepository
	Article   repository.ArticleRepository
	Lineup    repository.LineupRepository
	Match     repository.MatchRepository
	TeamStats repository.TeamStatsRepository
}

type Services struct {
	JWT       *jwt.JWTService
	Auth      service.AuthService
	User      service.UserService
	Role      service.RoleService
	Team      service.TeamService
	Player    service.PlayerService
	Season    service.SeasonService
	Article   service.ArticleService
	Lineup    service.LineupService
	Match     service.MatchService
	TeamStats service.TeamStatsService
}

type Handlers struct {
	User      *handler.UserHandler
	Role      *handler.RoleHandler
	Team      *handler.TeamHandler
	Player    *handler.PlayerHandler
	Season    *handler.SeasonHandler
	Article   *handler.ArticleHandler
	Lineup    *handler.LineupHandler
	Match     *handler.MatchHandler
	TeamStats *handler.TeamStatsHandler
}

// NewServer creates and configures a new server instance with middleware and security settings.
func NewServer() *Server {
	// Create a new Gin instance without any middleware
	r := gin.New()

	// Add recovery middleware
	r.Use(gin.Recovery())

	// Add our structured logger middleware
	r.Use(middleware.StructuredLogger())

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

// SetupRouter configures all API routes and their handlers.
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

// Start initializes the server with HTTP configuration and handles graceful shutdown.
func (s *Server) Start() {
	// Configurar las rutas antes de iniciar el servidor
	s.SetupRouter()

	// Configurar el servidor HTTP con tiempos de espera para seguridad
	server := createHTTPServer(s.Port, s.Router)

	// Iniciar el servidor en una goroutine
	go startServer(server)

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
		slog.Error("Failed to read JWT secret", "error", err)
		os.Exit(1)
	}
	if jwtSecret == nil {
		slog.Error("JWT secret is not defined")
		os.Exit(1)
	}
	return jwtSecret
}

func setupPprofServer() {
	if os.Getenv("GIN_MODE") != "release" {
		go func() {
			slog.Info("Starting pprof server", "address", ":6060")
			if err := http.ListenAndServe(":6060", nil); err != nil {
				slog.Error("Could not start pprof server", "error", err)
				os.Exit(1)
			}
		}()
	}
}

func getDBConnection() *gorm.DB {
	db, err := orm.GetDBInstance()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
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

func startServer(server *http.Server) {
	slog.Info("Server starting", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Could not start server", "error", err)
		os.Exit(1)
	}
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}

// Funciones para inicializar componentes

func initializeRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:      repository.NewUserRepository(db),
		Role:      repository.NewRoleRepository(db),
		Team:      repository.NewTeamRepository(db),
		Player:    repository.NewPlayerRepository(db),
		Season:    repository.NewSeasonRepository(db),
		Article:   repository.NewArticleRepository(db),
		Lineup:    repository.NewLineupRepository(db),
		Match:     repository.NewMatchRepository(db),
		TeamStats: repository.NewTeamStatsRepository(db),
	}
}

func initializeServices(repos *Repositories, jwtSecret []byte) *Services {
	jwtService := jwt.NewJWTService(string(jwtSecret))
	return &Services{
		JWT:       jwtService,
		Auth:      service.NewAuthService(repos.User, jwtService),
		User:      service.NewUserService(repos.User),
		Role:      service.NewRoleService(repos.Role),
		Team:      service.NewTeamService(repos.Team),
		Player:    service.NewPlayerService(repos.Player),
		Season:    service.NewSeasonService(repos.Season),
		Article:   service.NewArticleService(repos.Article),
		Lineup:    service.NewLineupService(repos.Lineup),
		Match:     service.NewMatchService(repos.Match),
		TeamStats: service.NewTeamStatsService(repos.TeamStats),
	}
}

func initializeHandlers(services *Services) *Handlers {
	return &Handlers{
		User:      handler.NewUserHandler(services.Auth, services.User, services.Role),
		Role:      handler.NewRoleHandler(services.Role),
		Team:      handler.NewTeamHandler(services.Team),
		Player:    handler.NewPlayerHandler(services.Player),
		Season:    handler.NewSeasonHandler(services.Season),
		Article:   handler.NewArticleHandler(services.Article),
		Lineup:    handler.NewLineupHandler(services.Lineup),
		Match:     handler.NewMatchHandler(services.Match),
		TeamStats: handler.NewTeamStatsHandler(services.TeamStats),
	}
}

func initializeRoutes(r *gin.Engine, handlers *Handlers) {
	router.InitializeUserRoutes(r, handlers.User)
	router.InitializeRoleRoutes(r, handlers.Role)
	router.InitializeTeamRoutes(r, handlers.Team)
	router.InitializePlayerRoutes(r, handlers.Player)
	router.InitializeSeasonRoutes(r, handlers.Season)
	router.InitializeArticleRoutes(r, handlers.Article)
	router.InitializeLineupRoutes(r, handlers.Lineup)
	router.InitializeMatchRoutes(r, handlers.Match)
	router.InitializeTeamStatsRoutes(r, handlers.TeamStats)
}

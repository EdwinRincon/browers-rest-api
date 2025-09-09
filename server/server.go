package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/middleware"
	router "github.com/EdwinRincon/browersfc-api/api/router"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/config"
	docs "github.com/EdwinRincon/browersfc-api/docs"
	"github.com/EdwinRincon/browersfc-api/domain"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
	"github.com/EdwinRincon/browersfc-api/pkg/orm"
)

type Server struct {
	Router    *gin.Engine
	Port      string
	JWTSecret []byte
}

// Repositories contains all data access layer dependencies (driven ports).
// Following hexagonal architecture, these represent the infrastructure adapters.
type Repositories struct {
	User           domain.UserRepository
	Role           domain.RoleRepository
	Team           domain.TeamRepository
	Player         domain.PlayerRepository
	PlayerTeam     domain.PlayerTeamRepository
	Season         domain.SeasonRepository
	Lineup         domain.LineupRepository
	Match          domain.MatchRepository
	Article        domain.ArticleRepository
	TeamStat       domain.TeamStatsRepository
	PlayerStat     domain.PlayerStatsRepository
	Authentication domain.AuthenticationRepository
}

// Services contains domain services (business rules) and auxiliary application services.
// Domain services implement core business logic, while application services handle cross-cutting concerns.
type Services struct {
	// Application services (cross-cutting concerns)
	JWT  *jwt.JWTService
	Auth service.AuthService
	// Domain services (core - business rules)
	AuthenticationDomain *domainservice.AuthenticationDomainService
	PlayerDomain         *domainservice.PlayerDomainService
	PlayerTeamDomain     *domainservice.PlayerTeamDomainService
	RoleDomain           *domainservice.RoleDomainService
	SeasonDomain         *domainservice.SeasonDomainService
	UserDomain           *domainservice.UserDomainService
	TeamDomain           *domainservice.TeamDomainService
	MatchDomain          *domainservice.MatchDomainService
	LineupDomain         *domainservice.LineupDomainService
	TeamStatDomain       *domainservice.TeamStatsDomainService
	PlayerStatDomain     *domainservice.PlayerStatsDomainService
	ArticleDomain        *domainservice.ArticleDomainService
}

// Handlers contains HTTP adapters (driving adapters).
// these represent the HTTP layer adapters.
type Handlers struct {
	User       *handler.UserHandler
	Role       *handler.RoleHandler
	Team       *handler.TeamHandler
	Player     *handler.PlayerHandler
	PlayerTeam *handler.PlayerTeamHandler
	Season     *handler.SeasonHandler
	Lineup     *handler.LineupHandler
	Match      *handler.MatchHandler
	TeamStat   *handler.TeamStatsHandler
	PlayerStat *handler.PlayerStatsHandler
	Article    *handler.ArticleHandler
}

// NewServer creates and configures a new server instance with middleware and security settings.
func NewServer() *Server {
	// Create a new Gin instance without any middleware
	r := gin.New()

	// Add recovery middleware
	r.Use(gin.Recovery())

	// Add our structured logger middleware
	r.Use(middleware.StructuredLogger())

	// Add CORS and security headers middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())

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
	initializeRoutes(s.Router, handlers, services)

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

// =====================================================
// Dependency Initialization Functions
// =====================================================

// initializeRepositories creates and configures all repository instances.
// This represents the driven ports (infrastructure adapters)
func initializeRepositories(db *gorm.DB) *Repositories {
	// Create role repository first as it's needed by authentication repository
	roleRepo := persistence.NewRoleRepository(db)

	return &Repositories{
		User:           persistence.NewUserRepository(db),
		Role:           roleRepo,
		Team:           persistence.NewTeamRepository(db),
		Player:         persistence.NewPlayerRepository(db),
		PlayerTeam:     persistence.NewPlayerTeamRepository(db),
		Season:         persistence.NewSeasonRepository(db),
		Article:        persistence.NewArticleRepository(db),
		Lineup:         persistence.NewLineupRepository(db),
		Match:          persistence.NewMatchRepository(db),
		TeamStat:       persistence.NewTeamStatsRepository(db),
		PlayerStat:     persistence.NewPlayerStatsRepository(db),
		Authentication: persistence.NewAuthenticationRepository(roleRepo),
	}
}

// initializeServices creates and configures domain services and supporting application services.
// Domain services implement core business logic, while application services handle cross-cutting concerns.
func initializeServices(repos *Repositories, jwtSecret []byte) *Services {
	jwtService := jwt.NewJWTService(string(jwtSecret))

	// Create domain services using domain factory (core business logic)
	roleDomainService := CreateRoleDomainService(repos.Role)
	seasonDomainService := CreateSeasonDomainService(repos.Season)
	userDomainService := CreateUserDomainService(repos.User)
	teamDomainService := CreateTeamDomainService(repos.Team)
	playerDomainService := CreatePlayerDomainService(repos.Player)
	playerTeamDomainService := CreatePlayerTeamDomainService(repos.PlayerTeam, repos.Player, repos.Team, repos.Season)
	lineupDomainService := CreateLineupDomainService(repos.Lineup, repos.Match, repos.Player)
	matchDomainService := CreateMatchDomainService(repos.Match)
	teamStatsDomainService := CreateTeamStatsDomainService(repos.TeamStat, repos.Team, repos.Season)
	playerStatsDomainService := CreatePlayerStatsDomainService(repos.PlayerStat, repos.Player, repos.Match, repos.Season, repos.Team)
	articleDomainService := CreateArticleDomainService(repos.Article, repos.Season)
	authenticationDomainService := CreateAuthenticationDomainService(repos.Authentication)

	return &Services{
		// Application services (cross-cutting concerns)
		JWT:  jwtService,
		Auth: service.NewAuthService(jwtService),
		// Domain services (core - business rules)
		AuthenticationDomain: authenticationDomainService,
		PlayerDomain:         playerDomainService,
		PlayerTeamDomain:     playerTeamDomainService,
		RoleDomain:           roleDomainService,
		SeasonDomain:         seasonDomainService,
		UserDomain:           userDomainService,
		TeamDomain:           teamDomainService,
		LineupDomain:         lineupDomainService,
		MatchDomain:          matchDomainService,
		TeamStatDomain:       teamStatsDomainService,
		PlayerStatDomain:     playerStatsDomainService,
		ArticleDomain:        articleDomainService,
	}
}

// initializeHandlers creates and configures all HTTP handler instances.
// This represents the driving adapters (HTTP layer)
func initializeHandlers(services *Services) *Handlers {
	return &Handlers{
		User:       handler.NewUserHandler(services.Auth, services.UserDomain, services.RoleDomain),
		Role:       handler.NewRoleHandler(services.RoleDomain),
		Team:       handler.NewTeamHandler(services.TeamDomain),
		Player:     handler.NewPlayerHandler(services.PlayerDomain),
		PlayerTeam: handler.NewPlayerTeamHandler(services.PlayerTeamDomain),
		Season:     handler.NewSeasonHandler(services.SeasonDomain),
		Lineup:     handler.NewLineupHandler(services.LineupDomain),
		Article:    handler.NewArticleHandler(services.ArticleDomain),
		Match:      handler.NewMatchHandler(services.MatchDomain),
		TeamStat:   handler.NewTeamStatsHandler(services.TeamStatDomain),
		PlayerStat: handler.NewPlayerStatsHandler(services.PlayerStatDomain),
	}
}

// initializeRoutes configures all API routes and middleware.
// This wires the driving adapters (handlers) to HTTP endpoints and applies infrastructure middleware.
func initializeRoutes(r *gin.Engine, handlers *Handlers, services *Services) {
	// Authentication middleware (infrastructure concern)
	authService := services.AuthenticationDomain

	// Route initialization (wiring handlers to endpoints)
	router.InitializeUserRoutes(r, handlers.User, authService)
	router.InitializeRoleRoutes(r, handlers.Role, authService)
	router.InitializeTeamRoutes(r, handlers.Team, authService)
	router.InitializePlayerRoutes(r, handlers.Player, authService)
	router.InitializePlayerTeamRoutes(r, handlers.PlayerTeam, authService)
	router.InitializeSeasonRoutes(r, handlers.Season, authService)
	router.InitializeLineupRoutes(r, handlers.Lineup, authService)
	router.InitializeArticleRoutes(r, handlers.Article, authService)
	router.InitializeMatchRoutes(r, handlers.Match, authService)
	router.InitializeTeamStatsRoutes(r, handlers.TeamStat, authService)
	router.InitializePlayerStatsRoutes(r, handlers.PlayerStat, authService)
}

// =====================================================
// Infrastructure Setup Helpers
// =====================================================

// setupSwagger configures Swagger documentation settings.
func setupSwagger() {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "BrowersFC API"
	docs.SwaggerInfo.Description = "API para la gestión de la liga de fútbol BrowersFC"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"https"}
}

// setupPprofServer starts a profiling server for development environments.
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

// getDBConnection establishes and returns a database connection.
func getDBConnection() *gorm.DB {
	db, err := orm.GetDBInstance()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	return db
}

// getJWTSecret reads and validates the JWT secret from configuration.
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

// =====================================================
// HTTP Lifecycle Utilities
// =====================================================

// createHTTPServer creates an HTTP server with security-focused timeouts.
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

// startServer starts the HTTP server and logs startup information.
func startServer(server *http.Server) {
	slog.Info("server_startup", "address", server.Addr, "event", "server_start")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server_error", "error", err, "event", "server_failure")
		os.Exit(1)
	}
}

// gracefulShutdown handles server shutdown with proper cleanup and timeout.
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

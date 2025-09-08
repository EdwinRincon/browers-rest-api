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
	router "github.com/EdwinRincon/browersfc-api/api/router"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/config"
	docs "github.com/EdwinRincon/browersfc-api/docs"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
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
	User       *persistence.UserRepositoryImpl
	Role       *persistence.RoleRepositoryImpl
	Team       *persistence.TeamRepositoryImpl
	Player     *persistence.PlayerRepositoryImpl
	PlayerTeam *persistence.PlayerTeamRepositoryImpl
	Season     *persistence.SeasonRepositoryImpl
	Article    persistence.ArticleRepository
	Lineup     persistence.LineupRepository
	Match      persistence.MatchRepository
	TeamStat   persistence.TeamStatsRepository
	PlayerStat persistence.PlayerStatsRepository
}

type Services struct {
	JWT  *jwt.JWTService
	Auth service.AuthService
	// Note: Following services removed until their entities are migrated to hexagonal architecture:
	// Player     service.PlayerService
	// PlayerTeam service.PlayerTeamService
	// Article    service.ArticleService
	// TeamStat   service.TeamStatsService
	// PlayerStat service.PlayerStatsService
	// Lineup     service.LineupService
	Match service.MatchService

	// Domain-based services (hexagonal architecture)
	PlayerDomain     *domainservice.PlayerDomainService
	PlayerTeamDomain *domainservice.PlayerTeamDomainService
	RoleDomain       *domainservice.RoleDomainService
	SeasonDomain     *domainservice.SeasonDomainService
	UserDomain       *domainservice.UserDomainService
	TeamDomain       *domainservice.TeamDomainService
}

type Handlers struct {
	User       *handler.UserHandler
	Role       *handler.RoleHandler
	Team       *handler.TeamHandler
	Player     *handler.PlayerHandler
	PlayerTeam *handler.PlayerTeamHandler
	Season     *handler.SeasonHandler
	// Article    *handler.ArticleHandler
	// Lineup     *handler.LineupHandler
	Match *handler.MatchHandler
	// TeamStat   *handler.TeamStatsHandler
	// PlayerStat *handler.PlayerStatsHandler
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
	slog.Info("server_startup", "address", server.Addr, "event", "server_start")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server_error", "error", err, "event", "server_failure")
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
		User:       persistence.NewUserRepository(db),
		Role:       persistence.NewRoleRepository(db),
		Team:       persistence.NewTeamRepository(db),
		Player:     persistence.NewPlayerRepository(db),
		PlayerTeam: persistence.NewPlayerTeamRepository(db),
		Season:     persistence.NewSeasonRepository(db),
		Article:    persistence.NewArticleRepository(db),
		Lineup:     persistence.NewLineupRepository(db),
		Match:      persistence.NewMatchRepository(db),
		TeamStat:   persistence.NewTeamStatsRepository(db),
		PlayerStat: persistence.NewPlayerStatsRepository(db),
	}
}

func initializeServices(repos *Repositories, jwtSecret []byte) *Services {
	jwtService := jwt.NewJWTService(string(jwtSecret))

	matchService := service.NewMatchService(repos.Match)

	// Create domain services using domain factory
	roleDomainService := CreateRoleDomainService(repos.Role)
	seasonDomainService := CreateSeasonDomainService(repos.Season)
	userDomainService := CreateUserDomainService(repos.User)
	teamDomainService := CreateTeamDomainService(repos.Team)
	playerDomainService := CreatePlayerDomainService(repos.Player)
	playerTeamDomainService := CreatePlayerTeamDomainService(repos.PlayerTeam, repos.Player, repos.Team, repos.Season)

	return &Services{
		JWT:  jwtService,
		Auth: service.NewAuthService(jwtService),
		// Note: Following services will fail compilation until their entities are migrated:
		// Player:     service.NewPlayerService(repos.Player, repos.PlayerTeam, repos.Season),
		// PlayerTeam: service.NewPlayerTeamService(repos.PlayerTeam, repos.Player, repos.Team, repos.Season),
		// Article:    service.NewArticleService(repos.Article, seasonService),
		Match: matchService,
		// TeamStat:   service.NewTeamStatsService(repos.TeamStat, repos.Team, repos.Season),
		// Lineup:     service.NewLineupService(repos.Lineup, matchService),
		// PlayerStat: service.NewPlayerStatsService(repos.PlayerStat, repos.Player, repos.Match, repos.Season, repos.Team),

		// Domain-based services (hexagonal architecture)
		PlayerDomain:     playerDomainService,
		PlayerTeamDomain: playerTeamDomainService,
		RoleDomain:       roleDomainService,
		SeasonDomain:     seasonDomainService,
		UserDomain:       userDomainService,
		TeamDomain:       teamDomainService,
	}
}

func initializeHandlers(services *Services) *Handlers {
	return &Handlers{
		User:       handler.NewUserHandler(services.Auth, services.UserDomain, services.RoleDomain),
		Role:       handler.NewRoleHandler(services.RoleDomain),
		Team:       handler.NewTeamHandler(services.TeamDomain),
		Player:     handler.NewPlayerHandler(services.PlayerDomain),
		PlayerTeam: handler.NewPlayerTeamHandler(services.PlayerTeamDomain),
		Season:     handler.NewSeasonHandler(services.SeasonDomain),
		// Article:    handler.NewArticleHandler(services.Article),
		// Lineup:     handler.NewLineupHandler(services.Lineup, services.Player, services.Match),
		Match: handler.NewMatchHandler(services.Match),
		// TeamStat:   handler.NewTeamStatsHandler(services.TeamStat),
		// PlayerStat: handler.NewPlayerStatsHandler(services.PlayerStat),
	}
}

func initializeRoutes(r *gin.Engine, handlers *Handlers) {
	router.InitializeUserRoutes(r, handlers.User)
	router.InitializeRoleRoutes(r, handlers.Role)
	router.InitializeTeamRoutes(r, handlers.Team)
	router.InitializePlayerRoutes(r, handlers.Player)
	router.InitializePlayerTeamRoutes(r, handlers.PlayerTeam)
	router.InitializeSeasonRoutes(r, handlers.Season)
	// router.InitializeArticleRoutes(r, handlers.Article)
	// router.InitializeLineupRoutes(r, handlers.Lineup)
	router.InitializeMatchRoutes(r, handlers.Match)
	// router.InitializeTeamStatsRoutes(r, handlers.TeamStat)
	// router.InitializePlayerStatsRoutes(r, handlers.PlayerStat)
}

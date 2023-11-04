package server

import (
	"log"
	"os"

	"github.com/EdwinRincon/browersfc-api/api"
	"github.com/EdwinRincon/browersfc-api/api/handler"
	"github.com/EdwinRincon/browersfc-api/api/repository"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("No se ha definido la variable de entorno para los tokens")
	}

	return &Server{
		Router:    r,
		Port:      port,
		JWTSecret: jwtSecret,
	}
}

func (s *Server) Start() {
	s.Router.Run(":" + s.Port)
}

package service

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/repository"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService es la interfaz que define los métodos de autenticación
type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
	GenerateToken(username, role string) (string, error)
}

// AuthService es la implementación concreta de la interfaz AuthService
type authService struct {
	UserRepository repository.UserRepository
	JWTService     *jwt.JWTService
}

// NewAuthService crea una nueva instancia de AuthService
func NewAuthService(userRepo repository.UserRepository, jwtService *jwt.JWTService) AuthService {
	return &authService{
		UserRepository: userRepo,
		JWTService:     jwtService,
	}
}

// Authenticate verifica las credenciales del usuario y devuelve un token JWT si son válidas
func (s *authService) Authenticate(ctx context.Context, username, password string) (string, error) {
	// Obtener el usuario almacenado
	storedUser, err := s.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	// Verificar si el usuario fue encontrado
	if storedUser == nil {
		return "", errors.New("user not found")
	}

	// Verificar si la cuenta está bloqueada
	if storedUser.FailedLoginAttempts >= 5 {
		return "", errors.New("account locked due to too many failed login attempts")
	}

	// Comparar la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(password)); err != nil {
		// Incrementar el número de intentos fallidos
		storedUser.FailedLoginAttempts++
		_, updateErr := s.UserRepository.UpdateUser(ctx, storedUser)
		if updateErr != nil {
			return "", errors.New("failed to update user attempts")
		}
		return "", errors.New("invalid username or password")
	}

	// Restablecer el número de intentos fallidos en caso de éxito
	storedUser.FailedLoginAttempts = 0
	_, updateErr := s.UserRepository.UpdateLoginAttemps(ctx, storedUser)
	if updateErr != nil {
		return "", errors.New("failed to reset user attempts")
	}

	// Generar el token
	token, err := s.JWTService.GenerateToken(username, storedUser.Roles.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GenerateToken genera un token JWT para un usuario autenticado
func (s *authService) GenerateToken(username, role string) (string, error) {
	return s.JWTService.GenerateToken(username, role)
}

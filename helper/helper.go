package helper

import (
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResponseJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseJSONSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(200, ResponseJSON{message, data})
}

func IsStrongPassword(password string) bool {
	var hasUpperCase = false
	var hasLowerCase = false
	var hasDigit = false
	var hasSpecialChar = false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		}
		if unicode.IsLower(char) {
			hasLowerCase = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecialChar = true
		}
		if len(password) < 8 {
			return false
		}
	}
	if !hasUpperCase || !hasLowerCase || !hasDigit || !hasSpecialChar {
		return false
	}
	return true
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

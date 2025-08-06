package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type ResponseJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Items      any   `json:"items"`
	TotalCount int64 `json:"total_count"`
}

// TODO: May remove HashPassword and CheckPasswordHash, GenerateRandomState, we are using OAuth
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

// GenerateRandomState generates a cryptographically secure random state for OAuth
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		slog.Error("Failed to generate random state", "error", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateRandomPassword() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		slog.Error("Failed to generate random password", "error", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// getDBColumnMap extracts database column names from a GORM model struct
// Returns a map where keys are valid column names and values are true
func getDBColumnMap(model any) map[string]bool {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make(map[string]bool)
	parseFields(t, fields)
	return fields
}

func parseFields(t reflect.Type, fields map[string]bool) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Recurse into embedded structs
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			parseFields(field.Type, fields)
			continue
		}

		tag := field.Tag.Get("gorm")
		columnName := ""

		// Extract column name from gorm tag
		for _, part := range strings.Split(tag, ";") {
			if strings.HasPrefix(part, "column:") {
				columnName = strings.TrimPrefix(part, "column:")
				break
			}
		}

		// If no explicit column name, use toSnakeCase of field name
		if columnName == "" {
			columnName = toSnakeCase(field.Name)
		}

		fields[columnName] = true
	}
}

// toSnakeCase converts a camelCase or PascalCase string to snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// ValidateSort checks if the provided sort field is valid for the given model
// Returns an error if the sort field is invalid
func ValidateSort(model any, sort string) error {
	validSorts := getDBColumnMap(model)
	if !validSorts[sort] {
		return fmt.Errorf("invalid sort field: %s", sort)
	}
	return nil
}

package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"unicode"
)

const (
	gormColumnPrefix = "column:"
)

type ResponseJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Items      any   `json:"items"`
	TotalCount int64 `json:"total_count"`
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

// GenerateRandomState generates a cryptographically secure random state for OAuth
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		slog.Error("Failed to generate random state", "error", err)
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

// parseFields recursively parses struct fields and extracts valid database columns
func parseFields(t reflect.Type, fields map[string]bool) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if shouldSkipField(field) {
			continue
		}

		if shouldProcessEmbeddedStruct(field) {
			parseFields(field.Type, fields)
			continue
		}

		columnName := extractColumnName(field)
		if columnName != "" {
			fields[columnName] = true
		}
	}
}

// shouldSkipField determines if a field should be skipped during parsing
func shouldSkipField(field reflect.StructField) bool {
	return isUnexportedField(field) || isGormRelationshipField(field) || isInvalidStructField(field)
}

// isUnexportedField checks if the field is unexported
func isUnexportedField(field reflect.StructField) bool {
	return field.PkgPath != ""
}

// shouldProcessEmbeddedStruct checks if field is an embedded struct that should be recursively processed
func shouldProcessEmbeddedStruct(field reflect.StructField) bool {
	return field.Anonymous && field.Type.Kind() == reflect.Struct
}

// isGormRelationshipField checks if the field is a GORM relationship field
func isGormRelationshipField(field reflect.StructField) bool {
	tag := field.Tag.Get("gorm")
	relationshipTags := []string{"foreignKey:", "references:", "many2many:", "constraint:"}

	if tag == "-" {
		return true
	}

	for _, relTag := range relationshipTags {
		if strings.Contains(tag, relTag) {
			return true
		}
	}
	return false
}

// isInvalidStructField checks if a struct field should be excluded from column mapping
func isInvalidStructField(field reflect.StructField) bool {
	fieldType := getActualFieldType(field.Type)

	if fieldType.Kind() != reflect.Struct {
		return false
	}

	tag := field.Tag.Get("gorm")
	return !isValidStructField(fieldType, tag)
}

// getActualFieldType dereferences pointer types to get the actual type
func getActualFieldType(fieldType reflect.Type) reflect.Type {
	if fieldType.Kind() == reflect.Ptr {
		return fieldType.Elem()
	}
	return fieldType
}

// isValidStructField determines if a struct field is valid for column mapping
func isValidStructField(fieldType reflect.Type, tag string) bool {
	isTimeType := fieldType.Name() == "Time" && fieldType.PkgPath() == "time"
	hasColumnTag := strings.Contains(tag, gormColumnPrefix)
	isEmbedded := strings.Contains(tag, "embedded")

	return isTimeType || hasColumnTag || isEmbedded
}

// extractColumnName extracts the database column name from a field
func extractColumnName(field reflect.StructField) string {
	tag := field.Tag.Get("gorm")

	// Extract column name from gorm tag
	if columnName := extractColumnNameFromTag(tag); columnName != "" {
		return columnName
	}

	// Use snake_case of field name as fallback
	return toSnakeCase(field.Name)
}

// extractColumnNameFromTag extracts the column name from a GORM tag
func extractColumnNameFromTag(tag string) string {
	for _, part := range strings.Split(tag, ";") {
		if strings.HasPrefix(part, gormColumnPrefix) {
			return strings.TrimPrefix(part, gormColumnPrefix)
		}
	}
	return ""
}

// toSnakeCase converts a camelCase or PascalCase string to snake_case
func toSnakeCase(str string) string {
	var result []rune
	runes := []rune(str)

	for i, r := range runes {
		if shouldAddUnderscore(runes, i) {
			result = append(result, '_')
		}

		if unicode.IsUpper(r) {
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// shouldAddUnderscore determines if an underscore should be added before the current character
func shouldAddUnderscore(runes []rune, i int) bool {
	if i == 0 {
		return false
	}

	curr := runes[i]
	prev := runes[i-1]

	// Number after letter or letter after number
	if (unicode.IsNumber(curr) && unicode.IsLetter(prev)) ||
		(unicode.IsLetter(curr) && unicode.IsNumber(prev)) {
		return true
	}

	// Uppercase rules
	if unicode.IsUpper(curr) {
		// Previous is lowercase
		if unicode.IsLower(prev) {
			return true
		}
		// Previous is uppercase and next is lowercase
		if i+1 < len(runes) && unicode.IsUpper(prev) && unicode.IsLower(runes[i+1]) {
			return true
		}
	}

	return false
}

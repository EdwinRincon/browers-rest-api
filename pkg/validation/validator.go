package validation

import (
	"net/mail"
	"reflect"
	"regexp"
	"strings"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/go-playground/validator/v10"
)

var (
	validate   *validator.Validate
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	urlRegex   = regexp.MustCompile(`^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)
)

// InitValidator initializes the validator with custom validations
func InitValidator() {
	validate = validator.New()

	// Register custom validations
	_ = validate.RegisterValidation("safe_email", validateSafeEmail)
	_ = validate.RegisterValidation("allowed_domain", validateAllowedDomain)
	_ = validate.RegisterValidation("safe_url", validateSafeURL)

	// Register custom tag name functions
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	if validate == nil {
		InitValidator()
	}
	return validate.Struct(s)
}

// validateSafeEmail ensures email follows security best practices
func validateSafeEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// Check basic format
	if !emailRegex.MatchString(email) {
		return false
	}

	// Parse email to validate structure
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return false
	}

	// Check email length
	if len(email) > 254 { // RFC 5321
		return false
	}

	// Check for potentially dangerous characters
	if strings.ContainsAny(email, "<>()[]\\,;:") {
		return false
	}

	localPart := strings.Split(email, "@")[0]
	if len(localPart) > 64 { // RFC 5321
		return false
	}

	return true
}

// validateAllowedDomain checks if email domain is in allowed list
func validateAllowedDomain(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := strings.ToLower(parts[1])
	for _, allowed := range constants.AllowedEmailDomains {
		if domain == allowed {
			return true
		}
	}

	return false
}

// validateSafeURL validates URLs
func validateSafeURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	if url == "" {
		return true // Allow empty values
	}

	return urlRegex.MatchString(url)
}

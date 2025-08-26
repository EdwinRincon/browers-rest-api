package validation

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/go-playground/validator/v10"
)

var (
	validate     *validator.Validate
	validateOnce sync.Once
	validateErr  error
	emailRegex   = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	urlRegex     = regexp.MustCompile(`^https?://(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}(?:/[a-zA-Z0-9\-._~!$&'()*+,;=:@%/?]*)?$`)
)

// InitValidator initializes the validator with custom validations
func InitValidator() error {
	validate = validator.New()

	// Register custom validations
	if err := validate.RegisterValidation("safe_email", validateSafeEmail); err != nil {
		return err
	}
	if err := validate.RegisterValidation("allowed_domain", validateAllowedDomain); err != nil {
		return err
	}
	if err := validate.RegisterValidation("safe_url", validateSafeURL); err != nil {
		return err
	}

	// Register custom tag name functions
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return nil
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	validateOnce.Do(func() {
		validateErr = InitValidator()
	})

	if validateErr != nil {
		return validateErr
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
	if domain == "" {
		return false
	}
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

// ExtractValidationErrors converts go-playground/validator errors into a user-friendly map
func ExtractValidationErrors(err error) map[string]string {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	errorMap := make(map[string]string)
	for _, fieldErr := range validationErrors {
		fieldName := getFieldName(fieldErr)
		message := getErrorMessage(fieldErr)
		errorMap[fieldName] = message
	}

	return errorMap
}

// getFieldName extracts the field name from a validation error
func getFieldName(fieldErr validator.FieldError) string {
	fieldName := fieldErr.Field()

	structType, ok := underlyingStruct(fieldErr.Value())
	if !ok {
		return fieldName
	}

	if field, found := structType.FieldByName(fieldErr.StructField()); found {
		if jsonTag := jsonTagName(field.Tag.Get("json")); jsonTag != "" {
			fieldName = jsonTag
		}
	}

	return fieldName
}

// underlyingStruct returns the underlying struct type and a boolean indicating success
func underlyingStruct(value interface{}) (reflect.Type, bool) {
	structType := reflect.TypeOf(value)
	if structType == nil {
		return nil, false
	}
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	return structType, structType.Kind() == reflect.Struct
}

func jsonTagName(tag string) string {
	if tag == "" || tag == "-" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}

// messageFunc defines a function type for generating error messages
type messageFunc func(fieldErr validator.FieldError) string

// registry of validation tag -> message generator
var messageRegistry = map[string]messageFunc{
	"required": func(_ validator.FieldError) string {
		return "This field is required"
	},
	"min": func(fe validator.FieldError) string {
		return fmt.Sprintf("Should be at least %s characters long", fe.Param())
	},
	"max": func(fe validator.FieldError) string {
		return fmt.Sprintf("Should not exceed %s characters", fe.Param())
	},
	"email": func(_ validator.FieldError) string {
		return "Invalid email format"
	},
	"url": func(_ validator.FieldError) string {
		return "Invalid URL format"
	},
	"safe_url": func(_ validator.FieldError) string {
		return "Invalid or unsafe URL"
	},
	"allowed_domain": func(_ validator.FieldError) string {
		return "Email domain is not allowed"
	},
	"safe_email": func(_ validator.FieldError) string {
		return "Invalid or unsafe email format"
	},
	"gte": func(fe validator.FieldError) string {
		return fmt.Sprintf("Should be greater than or equal to %s", fe.Param())
	},
	"lte": func(fe validator.FieldError) string {
		return fmt.Sprintf("Should be less than or equal to %s", fe.Param())
	},
	"alphanum": func(_ validator.FieldError) string {
		return "Should contain only alphanumeric characters"
	},
	"len": func(fe validator.FieldError) string {
		return fmt.Sprintf("Should be exactly %s characters long", fe.Param())
	},
	"oneof": func(fe validator.FieldError) string {
		return fmt.Sprintf("Should be one of: %s", fe.Param())
	},
}

func getErrorMessage(fieldErr validator.FieldError) string {
	if msgFunc, exists := messageRegistry[fieldErr.Tag()]; exists {
		return msgFunc(fieldErr)
	}
	// fallback if tag is not in registry
	return fmt.Sprintf("Failed validation on '%s'", fieldErr.Tag())
}

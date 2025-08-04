package constants

// OAuth2 Security Constants
const (
	// Email domain restrictions for account creation
	MaxLoginAttemptsPerIP  = 5  // Maximum login attempts per IP address
	LoginLockoutDuration   = 15 // Lockout duration in minutes
	MinPasswordLength      = 8  // Minimum password length
	TokenExpiryHours       = 24 // JWT token expiry in hours
	MaxSessionsPerUser     = 5  // Maximum concurrent sessions per user
	OAuthStateDuration     = 10 // OAuth state validity duration in minutes
	MaxNewAccountsPerIP    = 2  // Maximum new accounts per IP per day
	MaxNewAccountsPerEmail = 1  // Maximum accounts per email domain
)

// Allowed email domains for account creation
var AllowedEmailDomains = []string{
	"gmail.com",
	"outlook.com",
	"hotmail.com",
	"yahoo.com",
}

// Security Headers
var SecurityHeaders = map[string]string{
	"Content-Security-Policy":           "default-src 'self'; img-src 'self' https://*.googleusercontent.com; script-src 'self'",
	"X-Frame-Options":                   "DENY",
	"X-Content-Type-Options":            "nosniff",
	"Referrer-Policy":                   "strict-origin-when-cross-origin",
	"X-XSS-Protection":                  "1; mode=block",
	"X-Permitted-Cross-Domain-Policies": "none",
	"Access-Control-Allow-Origin":       "*", //TODO: Restrict this in production
}

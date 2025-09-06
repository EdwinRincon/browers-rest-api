package security

// OAuth2 Security Constants
const (
	// Security limits and durations
	MaxLoginAttemptsPerIP  = 5  // Maximum login attempts per IP address
	LoginLockoutDuration   = 15 // Lockout duration in minutes
	MinPasswordLength      = 8  // Minimum password length
	TokenExpiryHours       = 24 // JWT token expiry in hours
	MaxSessionsPerUser     = 5  // Maximum concurrent sessions per user
	OAuthStateDuration     = 10 // OAuth state validity duration in minutes
	MaxNewAccountsPerIP    = 2  // Maximum new accounts per IP per day
	MaxNewAccountsPerEmail = 1  // Maximum accounts per email domain
)

// AllowedEmailDomains defines the email domains allowed for account creation
var AllowedEmailDomains = []string{
	"gmail.com",
	"outlook.com",
	"hotmail.com",
	"yahoo.com",
}

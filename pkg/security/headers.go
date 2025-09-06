package security

// SecurityHeaders defines default security headers for HTTP responses
var SecurityHeaders = map[string]string{
	"Content-Security-Policy":           "default-src 'self'; img-src 'self' https://*.googleusercontent.com; script-src 'self'",
	"X-Frame-Options":                   "DENY",
	"X-Content-Type-Options":            "nosniff",
	"Referrer-Policy":                   "strict-origin-when-cross-origin",
	"X-XSS-Protection":                  "1; mode=block",
	"X-Permitted-Cross-Domain-Policies": "none",
	"Access-Control-Allow-Origin":       "*", //TODO: Restrict this in production
}

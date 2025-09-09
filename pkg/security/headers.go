package security

// SecurityHeaders defines default security headers for HTTP responses
var SecurityHeaders = map[string]string{
	"Content-Security-Policy":           "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https:; connect-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'",
	"X-Frame-Options":                   "DENY",
	"X-Content-Type-Options":            "nosniff",
	"Referrer-Policy":                   "strict-origin-when-cross-origin",
	"X-Permitted-Cross-Domain-Policies": "none",
	"Access-Control-Allow-Origin":       "*", //TODO: Restrict this in production
}

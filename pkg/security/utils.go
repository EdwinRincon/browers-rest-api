package security

import (
	"slices"
	"strings"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/gin-gonic/gin"
)

// ValidateEmailDomain checks if the email domain is allowed for account creation
func ValidateEmailDomain(email string) bool {
	_, domain := splitEmail(email)
	if domain == "" {
		return false
	}

	return slices.Contains(AllowedEmailDomains, domain)
}

// SetSecureCookie sets a cookie with security configurations from the app config
func SetSecureCookie(c *gin.Context, name, value string, maxAge int) {
	cfg := config.Config.CookieConfig
	c.SetCookie(
		name,
		value,
		maxAge,
		cfg.Path,
		cfg.Domain,
		cfg.Secure,
		cfg.HTTPOnly,
	)
}

// splitEmail splits an email address into local part and domain
func splitEmail(email string) (string, string) {
	username, domain, found := strings.Cut(email, "@")
	if !found || username == "" || domain == "" {
		return "", ""
	}
	return username, domain
}

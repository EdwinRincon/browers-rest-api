package auth

import (
	"slices"
	"strings"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/gin-gonic/gin"
)

// Security package provides authentication and security related utilities

// ValidateEmailDomain checks if the email domain is allowed
func ValidateEmailDomain(email string) bool {
	_, domain := splitEmail(email)
	if domain == "" {
		return false
	}

	return slices.Contains(constants.AllowedEmailDomains, domain)
}

// SetSecureCookie sets a cookie with security configurations
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

// splitEmail splits an email into local part and domain
func splitEmail(email string) (string, string) {
	username, domain, found := strings.Cut(email, "@")
	if !found || username == "" || domain == "" {
		return "", ""
	}
	return username, domain
}

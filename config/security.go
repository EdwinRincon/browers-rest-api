package config

import (
	"os"
	"strings"
	"time"
)

type SecurityConfig struct {
	IsDevelopment bool
	CSPConfig     ContentSecurityPolicyConfig
	CookieConfig  CookieSecurityConfig
}

type ContentSecurityPolicyConfig struct {
	DefaultSrc     []string
	ScriptSrc      []string
	StyleSrc       []string
	ImgSrc         []string
	ConnectSrc     []string
	FontSrc        []string
	ObjectSrc      []string
	MediaSrc       []string
	FrameSrc       []string
	FrameAncestors []string
}

type CookieSecurityConfig struct {
	Secure   bool
	HTTPOnly bool
	SameSite string
	Path     string
	Domain   string
	MaxAge   time.Duration
}

var Config SecurityConfig

func InitSecurityConfig() {
	isDev := os.Getenv("GIN_MODE") != "release"

	Config = SecurityConfig{
		IsDevelopment: isDev,
		CSPConfig: ContentSecurityPolicyConfig{
			DefaultSrc:     []string{"'self'"},
			ScriptSrc:      []string{"'self'", "'unsafe-inline'", "https://apis.google.com"},
			StyleSrc:       []string{"'self'", "'unsafe-inline'"},
			ImgSrc:         []string{"'self'", "https://*.googleusercontent.com", "data:"},
			ConnectSrc:     []string{"'self'", "https://accounts.google.com"},
			FontSrc:        []string{"'self'"},
			ObjectSrc:      []string{"'none'"},
			MediaSrc:       []string{"'self'"},
			FrameSrc:       []string{"'self'", "https://accounts.google.com"},
			FrameAncestors: []string{"'none'"},
		},
		CookieConfig: CookieSecurityConfig{
			Secure:   !isDev,
			HTTPOnly: true,
			SameSite: "Lax",
			Path:     "/",
			Domain:   "",
			MaxAge:   24 * time.Hour,
		},
	}
}

// BuildCSPHeader constructs the Content-Security-Policy header value
// using the current configuration. It returns a string suitable for use
// as the value of the CSP HTTP header.
func BuildCSPHeader() string {
	c := Config.CSPConfig

	// Collect directives
	directives := []struct {
		name    string
		sources []string
	}{
		{"default-src", c.DefaultSrc},
		{"script-src", c.ScriptSrc},
		{"style-src", c.StyleSrc},
		{"img-src", c.ImgSrc},
		{"connect-src", c.ConnectSrc},
		{"font-src", c.FontSrc},
		{"object-src", c.ObjectSrc},
		{"media-src", c.MediaSrc},
		{"frame-src", c.FrameSrc},
		{"frame-ancestors", c.FrameAncestors},
	}

	var parts []string
	for _, d := range directives {
		if len(d.sources) > 0 {
			// Only include if sources is not empty
			parts = append(parts, d.name+" "+strings.Join(d.sources, " "))
		}
	}
	return strings.Join(parts, "; ")
}

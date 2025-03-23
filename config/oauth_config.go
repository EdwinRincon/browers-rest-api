package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	OAuthConfig *oauth2.Config
	pkceStore   sync.Map // Thread-safe map for storing PKCE verifiers
)

type PKCEParams struct {
	Verifier  string
	Challenge string
}

func InitOAuth() error {
	clientID := os.Getenv("OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return errors.New("missing required OAuth configuration")
	}

	OAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	return nil
}

func GeneratePKCE() (*PKCEParams, error) {
	verifier, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return &PKCEParams{
		Verifier:  verifier,
		Challenge: challenge,
	}, nil
}

func StorePKCE(state string, params *PKCEParams) {
	pkceStore.Store(state, params)
}

func GetAndDeletePKCE(state string) (*PKCEParams, bool) {
	if value, ok := pkceStore.LoadAndDelete(state); ok {
		return value.(*PKCEParams), true
	}
	return nil, false
}

// generateRandomString creates a random string of specified length
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}

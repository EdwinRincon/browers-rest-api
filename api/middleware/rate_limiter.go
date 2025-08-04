// Package middleware provides HTTP middleware functions for the BrowersFC API
package middleware

import (
	"log/slog"
	"sync"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// limiterEntry tracks both the rate limiter and its last activity time
type limiterEntry struct {
	limiter      *rate.Limiter
	lastActivity time.Time
}

// IPRateLimiter implements a thread-safe rate limiter for IP addresses
// using the token bucket algorithm. It automatically cleans up old entries
// to prevent memory leaks.
type IPRateLimiter struct {
	sync.RWMutex
	limiters map[string]*limiterEntry
	r        rate.Limit    // requests per second
	b        int           // burst size
	ttl      time.Duration // time-to-live for limiter entries
	metrics  *RateLimitMetrics
}

// rateLimitConfig holds configuration for a rate limiter middleware
type rateLimitConfig struct {
	limiter *IPRateLimiter
	message string
}

// Allow checks if a request should be allowed and updates metrics
func (entry *limiterEntry) Allow() bool {
	if allowed := entry.limiter.Allow(); allowed {
		entry.lastActivity = time.Now()
		return true
	}
	return false
}

// isExpired checks if the entry has exceeded its TTL
func (entry *limiterEntry) isExpired(ttl time.Duration) bool {
	return time.Since(entry.lastActivity) > ttl
}

// RateLimitMetrics tracks rate limiting statistics for monitoring
type RateLimitMetrics struct {
	sync.RWMutex
	TotalRequests   uint64
	BlockedRequests uint64
	ActiveLimiters  uint64
}

// NewIPRateLimiter creates a new rate limiter with specified rate, burst size and TTL
func NewIPRateLimiter(r rate.Limit, b int, ttl time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters: make(map[string]*limiterEntry),
		r:        r,
		b:        b,
		ttl:      ttl,
		metrics:  &RateLimitMetrics{},
	}

	// Start background cleanup
	go limiter.cleanup()

	return limiter
}

// getLimiter returns the rate limiter for the specified IP address
func (i *IPRateLimiter) getLimiter(ip string) *limiterEntry {
	i.RLock()
	entry, exists := i.limiters[ip]
	i.RUnlock()

	if exists {
		return entry
	}

	i.Lock()
	defer i.Unlock()

	// Double-check after acquiring write lock
	entry, exists = i.limiters[ip]
	if exists {
		return entry
	}

	// Create new limiter entry
	entry = &limiterEntry{
		limiter:      rate.NewLimiter(i.r, i.b),
		lastActivity: time.Now(),
	}
	i.limiters[ip] = entry

	// Update metrics
	i.metrics.Lock()
	i.metrics.ActiveLimiters++
	i.metrics.Unlock()

	slog.Debug("new rate limiter created",
		"ip", ip,
		"rate", i.r,
		"burst", i.b,
		"active_limiters", i.metrics.ActiveLimiters)

	return entry
}

// cleanup periodically removes expired limiters to prevent memory leaks
func (i *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(i.ttl)
	defer ticker.Stop()

	for range ticker.C {
		expired := make([]string, 0)

		i.RLock()
		for ip, entry := range i.limiters {
			if entry.isExpired(i.ttl) {
				expired = append(expired, ip)
			}
		}
		i.RUnlock()

		if len(expired) > 0 {
			i.Lock()
			for _, ip := range expired {
				delete(i.limiters, ip)
			}
			i.Unlock()

			// Update metrics
			i.metrics.Lock()
			i.metrics.ActiveLimiters -= uint64(len(expired))
			i.metrics.Unlock()

			slog.Debug("cleaned up expired rate limiters",
				"count", len(expired),
				"active_limiters", i.metrics.ActiveLimiters)
		}
	}
}

// GetMetrics returns the current rate limiting metrics
func (i *IPRateLimiter) GetMetrics() RateLimitMetrics {
	i.metrics.RLock()
	defer i.metrics.RUnlock()
	return RateLimitMetrics{
		TotalRequests:   i.metrics.TotalRequests,
		BlockedRequests: i.metrics.BlockedRequests,
		ActiveLimiters:  i.metrics.ActiveLimiters,
	}
}

// makeRateLimitHandler creates a middleware handler with the specified config
func makeRateLimitHandler(config rateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		entry := config.limiter.getLimiter(ip)

		// Update metrics before checking allow
		config.limiter.metrics.Lock()
		config.limiter.metrics.TotalRequests++
		config.limiter.metrics.Unlock()

		if !entry.Allow() {
			// Update blocked requests metric
			config.limiter.metrics.Lock()
			config.limiter.metrics.BlockedRequests++
			config.limiter.metrics.Unlock()

			c.AbortWithStatusJSON(429, gin.H{
				"error": config.message,
			})
			return
		}
		c.Next()
	}
}

// Global rate limiters
var (
	// Generic rate limiter: 100 requests per minute per IP
	globalLimiter = NewIPRateLimiter(rate.Every(time.Minute), 100, 1*time.Hour)
	// Specific rate limiters for sensitive operations
	authLimiter       = NewIPRateLimiter(rate.Every(time.Minute), constants.MaxLoginAttemptsPerIP, 15*time.Minute)
	newAccountLimiter = NewIPRateLimiter(rate.Every(24*time.Hour), constants.MaxNewAccountsPerIP, 24*time.Hour)
)

// RateLimit is a generic rate limiter for all API endpoints
func RateLimit() gin.HandlerFunc {
	return makeRateLimitHandler(rateLimitConfig{
		limiter: globalLimiter,
		message: "Too many requests, please try again later",
	})
}

// RateLimitAuth handles rate limiting for authentication requests
func RateLimitAuth() gin.HandlerFunc {
	return makeRateLimitHandler(rateLimitConfig{
		limiter: authLimiter,
		message: "Too many authentication attempts, please try again later",
	})
}

// RateLimitNewAccounts handles rate limiting for new account creation
func RateLimitNewAccounts() gin.HandlerFunc {
	return makeRateLimitHandler(rateLimitConfig{
		limiter: newAccountLimiter,
		message: "Maximum number of new accounts reached for today",
	})
}

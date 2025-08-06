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

// limiterEntry holds the rate limiter and last activity timestamp
type limiterEntry struct {
	limiter      *rate.Limiter
	lastActivity time.Time
}

// IPRateLimiter manages rate limiting per IP
type IPRateLimiter struct {
	sync.RWMutex
	limiters map[string]*limiterEntry
	r        rate.Limit
	b        int
	ttl      time.Duration
	metrics  RateLimitMetrics
}

// RateLimitMetrics tracks limiter statistics
type RateLimitMetrics struct {
	sync.RWMutex
	TotalRequests   uint64
	BlockedRequests uint64
	ActiveLimiters  uint64
}

// NewIPRateLimiter constructs an IPRateLimiter
func NewIPRateLimiter(r rate.Limit, b int, ttl time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters: make(map[string]*limiterEntry),
		r:        r,
		b:        b,
		ttl:      ttl,
	}
	go limiter.cleanup()
	return limiter
}

// getLimiter fetches/creates a limiter for the IP
func (i *IPRateLimiter) getLimiter(ip string) *limiterEntry {
	i.RLock()
	entry := i.limiters[ip]
	i.RUnlock()
	if entry != nil {
		return entry
	}

	i.Lock()
	defer i.Unlock()
	if entry = i.limiters[ip]; entry == nil {
		entry = &limiterEntry{
			limiter:      rate.NewLimiter(i.r, i.b),
			lastActivity: time.Now(),
		}
		i.limiters[ip] = entry
		i.metrics.Lock()
		i.metrics.ActiveLimiters++
		i.metrics.Unlock()
		slog.Debug("Created new rate limiter", "ip", ip, "rate", i.r, "burst", i.b, "active_limiters", i.metrics.ActiveLimiters)
	}
	return entry
}

// Allow checks if the request is allowed and updates activity
func (e *limiterEntry) Allow() bool {
	if allowed := e.limiter.Allow(); allowed {
		e.lastActivity = time.Now()
		return true
	}
	return false
}

// isExpired returns true if the entry should be removed
func (e *limiterEntry) isExpired(ttl time.Duration) bool {
	return time.Since(e.lastActivity) > ttl
}

// cleanup periodically purges expired limiters
func (i *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(i.ttl)
	defer ticker.Stop()
	for range ticker.C {
		var expired []string
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
			i.metrics.Lock()
			i.metrics.ActiveLimiters -= uint64(len(expired))
			i.metrics.Unlock()
			i.Unlock()
			slog.Debug("Cleaned expired rate limiters", "count", len(expired), "active_limiters", i.metrics.ActiveLimiters)
		}
	}
}

// GetMetrics retrieves a snapshot of limiter metrics
func (i *IPRateLimiter) GetMetrics() RateLimitMetrics {
	i.metrics.RLock()
	defer i.metrics.RUnlock()
	return RateLimitMetrics{
		TotalRequests:   i.metrics.TotalRequests,
		BlockedRequests: i.metrics.BlockedRequests,
		ActiveLimiters:  i.metrics.ActiveLimiters,
	}
}

// rateLimitConfig configures the rate limit handler
type rateLimitConfig struct {
	limiter *IPRateLimiter
	message string
}

// makeRateLimitHandler returns a gin.HandlerFunc for rate limiting
func makeRateLimitHandler(config rateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		entry := config.limiter.getLimiter(ip)

		config.limiter.metrics.Lock()
		config.limiter.metrics.TotalRequests++
		config.limiter.metrics.Unlock()

		if !entry.Allow() {
			config.limiter.metrics.Lock()
			config.limiter.metrics.BlockedRequests++
			config.limiter.metrics.Unlock()
			c.AbortWithStatusJSON(429, gin.H{"error": config.message})
			return
		}
		c.Next()
	}
}

// Global rate limiters
var (
	globalLimiter     = NewIPRateLimiter(rate.Every(time.Minute), 100, time.Hour)
	authLimiter       = NewIPRateLimiter(rate.Every(time.Minute), constants.MaxLoginAttemptsPerIP, 15*time.Minute)
	newAccountLimiter = NewIPRateLimiter(rate.Every(24*time.Hour), constants.MaxNewAccountsPerIP, 24*time.Hour)
)

// RateLimit provides global rate limiting for all endpoints
func RateLimit() gin.HandlerFunc {
	return makeRateLimitHandler(rateLimitConfig{
		limiter: globalLimiter,
		message: "Too many requests, please try again later",
	})
}

// RateLimitAuth limits authentication attempts
func RateLimitAuth() gin.HandlerFunc {
	return makeRateLimitHandler(rateLimitConfig{
		limiter: authLimiter,
		message: "Too many authentication attempts, please try again later",
	})
}

// RateLimitNewAccounts limits account creation per IP
func RateLimitNewAccounts() gin.HandlerFunc {
	return makeRateLimitHandler(rateLimitConfig{
		limiter: newAccountLimiter,
		message: "Maximum number of new accounts reached for today",
	})
}

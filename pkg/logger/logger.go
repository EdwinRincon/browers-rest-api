package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

// LogConfig holds configuration for the application's logging
type LogConfig struct {
	Level   slog.Level
	Format  LogFormat
	Output  io.Writer
	IsDebug bool
}

// LogFormat represents the format of log output
type LogFormat string

const (
	JSONFormat LogFormat = "json"
	TextFormat LogFormat = "text"
)

// Context keys for request metadata
const (
	RequestIDKey  = "request_id"
	MethodKey     = "method"
	PathKey       = "path"
	StatusKey     = "status"
	DurationMsKey = "duration_ms"
	ClientIPKey   = "client_ip"
	UserAgentKey  = "user_agent"
	ErrorsKey     = "errors"

	// Rich error details key
	AppErrorKey = "app_error" // Key for storing rich error in Gin context
)

// LoggableError defines the interface for errors that can be logged with rich details
type LoggableError interface {
	// AddToLog adds structured error data to the given logger
	AddToLog(logger *slog.Logger) *slog.Logger
}

var defaultLogger *slog.Logger

// Setup initializes the global logger with the given configuration
func Setup(cfg LogConfig) {
	var handler slog.Handler

	// Create handler options with UTC time formatter
	handlerOptions := &slog.HandlerOptions{
		Level: cfg.Level,
		// Override the default time formatter to use UTC
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Convert time values to UTC before logging
			if a.Key == "time" && a.Value.Kind() == slog.KindTime {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.TimeValue(a.Value.Time().UTC()),
				}
			}
			return a
		},
	}

	// Configure handler based on format
	switch cfg.Format {
	case JSONFormat:
		handler = slog.NewJSONHandler(cfg.Output, handlerOptions)
	case TextFormat:
		handler = slog.NewTextHandler(cfg.Output, handlerOptions)
	default:
		// Default to JSON in production, text in development
		if cfg.IsDebug {
			handler = slog.NewTextHandler(cfg.Output, handlerOptions)
		} else {
			handler = slog.NewJSONHandler(cfg.Output, handlerOptions)
		}
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// GetLogger returns a logger with request context data if available
func GetLogger(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return defaultLogger
	}

	// Extract request ID from context
	if requestID, exists := GetRequestID(ctx); exists {
		return defaultLogger.With(slog.String(RequestIDKey, requestID))
	}

	return defaultLogger
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) (string, bool) {
	// Check if it's a Gin context
	if gc, ok := ctx.(*gin.Context); ok {
		if requestID, exists := gc.Get(RequestIDKey); exists {
			if id, ok := requestID.(string); ok {
				return id, true
			}
		}
		return "", false
	}

	// Check if it's in the context directly
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID, true
	}

	return "", false
}

// WithRequestID adds a request ID to the logger
func WithRequestID(l *slog.Logger, requestID string) *slog.Logger {
	return l.With(slog.String(RequestIDKey, requestID))
}

// WithHTTPRequest enriches logger with HTTP request metadata from Gin context
func WithHTTPRequest(l *slog.Logger, c *gin.Context, durationMs float64) *slog.Logger {
	// Build full path with query parameters if present
	path := c.Request.URL.Path
	if raw := c.Request.URL.RawQuery; raw != "" {
		path = path + "?" + raw
	}

	// Add HTTP request metadata
	logger := l.With(
		slog.String(MethodKey, c.Request.Method),
		slog.String(PathKey, path),
		slog.Int(StatusKey, c.Writer.Status()),
		slog.Float64(DurationMsKey, durationMs),
		slog.String(ClientIPKey, c.ClientIP()),
		slog.String(UserAgentKey, c.Request.UserAgent()),
	)

	// Add standard Gin errors if present
	if len(c.Errors) > 0 {
		logger = logger.With(slog.String(ErrorsKey, c.Errors.String()))
	}

	// Add rich application error details if present
	if appError, exists := c.Get(AppErrorKey); exists {
		if loggable, ok := appError.(LoggableError); ok {
			// Use the interface method to add error details to the logger
			logger = loggable.AddToLog(logger)
		}
	}

	return logger
}

// Debug logs at debug level with request ID if available
func Debug(ctx context.Context, msg string, args ...any) {
	GetLogger(ctx).Debug(msg, args...)
}

// Info logs at info level with request ID if available
func Info(ctx context.Context, msg string, args ...any) {
	GetLogger(ctx).Info(msg, args...)
}

// Warn logs at warning level with request ID if available
func Warn(ctx context.Context, msg string, args ...any) {
	GetLogger(ctx).Warn(msg, args...)
}

// Error logs at error level with request ID if available
func Error(ctx context.Context, msg string, args ...any) {
	GetLogger(ctx).Error(msg, args...)
}

// LogHTTPRequest logs an HTTP request with all its metadata
// This is the centralized function for logging HTTP requests
func LogHTTPRequest(c *gin.Context, durationMs float64, shouldSkip bool) {
	// Skip detailed logging for health checks and static files if successful
	if shouldSkip {
		return
	}

	// Get logger with request ID
	logger := GetLogger(c)

	// Add HTTP metadata
	httpLogger := WithHTTPRequest(logger, c, durationMs)

	// Determine message based on whether an error occurred
	var message string
	if _, exists := c.Get(AppErrorKey); exists {
		message = "http request error"
	} else if len(c.Errors) > 0 {
		message = "http request error"
	} else {
		message = "http request"
	}

	// Log with appropriate level based on status code
	statusCode := c.Writer.Status()
	switch {
	case statusCode >= 500:
		httpLogger.Error(message)
	case statusCode >= 400:
		httpLogger.Warn(message)
	default:
		httpLogger.Info(message)
	}
}

// StoreErrorForLogging stores an error in the Gin context for later logging
// This allows non-handler code to contribute error details without direct logging
func StoreErrorForLogging(c *gin.Context, err LoggableError) {
	c.Set(AppErrorKey, err)
}

// GetDefaultLevel returns the default log level based on environment
func GetDefaultLevel() slog.Level {
	if os.Getenv("GIN_MODE") == "release" {
		return slog.LevelInfo
	}
	return slog.LevelDebug
}

// GetDefaultFormat returns the default log format based on environment
func GetDefaultFormat() LogFormat {
	if os.Getenv("GIN_MODE") == "release" {
		return JSONFormat
	}
	return TextFormat
}

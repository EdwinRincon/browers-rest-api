package orm

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/middleware"
	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ContextAwareGormLogger implements the GORM logger interface with request context awareness.
// It adds request_id to logs when available and controls verbosity based on config.
type ContextAwareGormLogger struct {
	SlowThreshold time.Duration // Threshold for slow query logging in ms
	LogLevel      gormlogger.LogLevel
}

// NewContextAwareGormLogger creates a new GORM logger with contextual awareness.
func NewContextAwareGormLogger() *ContextAwareGormLogger {
	var level gormlogger.LogLevel

	// In production, only log slow queries and errors
	if config.AppLogConfig.Level > slog.LevelDebug {
		level = gormlogger.Error
	} else {
		level = gormlogger.Info
	}

	return &ContextAwareGormLogger{
		SlowThreshold: time.Duration(config.AppLogConfig.SlowQueryTime) * time.Millisecond,
		LogLevel:      level,
	}
}

// LogMode sets the log level for this logger
func (l *ContextAwareGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info messages
func (l *ContextAwareGormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		if ctx == nil {
			ctx = context.Background()
		}
		logger.Info(ctx, msg, args...)
	}
}

// Warn logs warning messages
func (l *ContextAwareGormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		if ctx == nil {
			ctx = context.Background()
		}
		logger.Warn(ctx, msg, args...)
	}
}

// Error logs error messages
func (l *ContextAwareGormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		if ctx == nil {
			ctx = context.Background()
		}
		logger.Error(ctx, msg, args...)
	}
}

// Trace logs SQL operations
func (l *ContextAwareGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Skip logging if below threshold
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// Get elapsed time
	elapsed := time.Since(begin)

	// Get SQL and rows affected
	sql, rows := fc()

	// Extract request ID from context if available
	var requestID string
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if id, exists := ginCtx.Get(middleware.RequestIDKey); exists {
			if idStr, ok := id.(string); ok {
				requestID = idStr
			}
		}
	}

	// Prepare log attributes
	attrs := []interface{}{
		"elapsed_ms", float64(elapsed.Nanoseconds()) / 1e6,
		"rows", rows,
	}

	if requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		// Log SQL errors (except "record not found" which is not really an error)
		logger.Error(ctx, "gorm query error", append(attrs,
			"sql", sql,
			"error", err.Error(),
		)...)
	case elapsed > l.SlowThreshold && l.SlowThreshold > 0 && l.LogLevel >= gormlogger.Warn:
		// Log slow queries as warnings
		logger.Warn(ctx, "gorm slow query", append(attrs,
			"sql", sql,
			"threshold_ms", float64(l.SlowThreshold.Nanoseconds())/1e6,
		)...)
	case l.LogLevel >= gormlogger.Info:
		// Log all other queries at debug level
		logger.Debug(ctx, "gorm query", append(attrs,
			"sql", sql,
		)...)
	}
}

// GetContextWithRequestID tries to extract a request ID from a gin context and add it
// to a new context for database operations.
func GetContextWithRequestID(ctx context.Context) context.Context {
	// If it's already a gin context, use it directly
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx
	}

	return ctx
}

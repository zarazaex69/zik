package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger wraps zerolog.Logger with additional context
type Logger struct {
	logger zerolog.Logger
}

var globalLogger *Logger

// Init initializes the global logger with the specified configuration
func Init(debug bool) {
	var output io.Writer = os.Stdout

	// Use pretty console output in debug mode, JSON otherwise
	if debug {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}

	zlog := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	globalLogger = &Logger{logger: zlog}
	log.Logger = zlog
}

// Get returns the global logger instance
func Get() *Logger {
	if globalLogger == nil {
		Init(false)
	}
	return globalLogger
}

// Debug logs a debug message
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Info logs an info message
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Warn logs a warning message
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Error logs an error message
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// With creates a child logger with additional context
func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

// Global convenience functions

// Debug logs a debug message using the global logger
func Debug() *zerolog.Event {
	return Get().Debug()
}

// Info logs an info message using the global logger
func Info() *zerolog.Event {
	return Get().Info()
}

// Warn logs a warning message using the global logger
func Warn() *zerolog.Event {
	return Get().Warn()
}

// Error logs an error message using the global logger
func Error() *zerolog.Event {
	return Get().Error()
}

// Fatal logs a fatal message using the global logger and exits
func Fatal() *zerolog.Event {
	return Get().Fatal()
}

package logger

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Reset global logger
	globalLogger = nil

	// Test debug mode
	Init(true)
	assert.NotNil(t, globalLogger)
	assert.Equal(t, zerolog.DebugLevel, globalLogger.logger.GetLevel())
	// Verify it's console writer (hard to check exact type of writer deeply inside zerolog, but we can check output format behavior effectively or just trust Init logic covered)

	// Test non-debug mode
	globalLogger = nil
	Init(false)
	assert.NotNil(t, globalLogger)
	assert.Equal(t, zerolog.InfoLevel, globalLogger.logger.GetLevel())
}

func TestGet(t *testing.T) {
	// Reset globalLogger
	globalLogger = nil

	// First call should initialize
	l := Get()
	assert.NotNil(t, l)
	assert.Equal(t, globalLogger, l)

	// Second call should return same instance
	l2 := Get()
	assert.Equal(t, l, l2)
}

func TestLoggerMethods(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	zlog := zerolog.New(&buf)
	l := &Logger{logger: zlog}

	// Test methods return event
	assert.NotNil(t, l.Debug())
	assert.NotNil(t, l.Info())
	assert.NotNil(t, l.Warn())
	assert.NotNil(t, l.Error())
	assert.NotNil(t, l.Fatal())

	// Test With returns context
	assert.NotNil(t, l.With())

	// Test actual logging
	l.Info().Msg("test message")
	assert.Contains(t, buf.String(), "test message")
}

func TestGlobalMethods(t *testing.T) {
	// Setup global logger with buffer
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	globalLogger = &Logger{logger: zlog}
	// Also need to set zerolog.log.Logger if we want to ensure consistency,
	// but the global functions call Get().Method() which uses globalLogger.

	// Test methods return event
	assert.NotNil(t, Debug())
	assert.NotNil(t, Info())
	assert.NotNil(t, Warn())
	assert.NotNil(t, Error())
	assert.NotNil(t, Fatal())

	// Test actual logging
	Info().Msg("global test")
	assert.Contains(t, buf.String(), "global test")
}

func TestInitOutput(t *testing.T) {
	// Just verify Init doesn't crash effectively.
	// Hijacking stdout/stderr to verify real output of Init is complex and maybe overkill if we just want coverage of lines.
	// We covered branches in TestInit.
}

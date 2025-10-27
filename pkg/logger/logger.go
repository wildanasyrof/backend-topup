package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// CtxKey is the key used to store the logger in the request context (e.g., Fiber's c.Locals)
const CtxKey = "logger"

type Fields map[string]any

type Logger interface {
	Info(msg string)
	Error(err error, msg string)
	Warn(msg string)
	Debug(msg string)
	Fatal(msg string)

	With(f Fields) Logger
	SetLevel(level string) // "debug", "info", "warn", "error"
}

// ZerologAdapter implements Logger.
type ZerologAdapter struct {
	logger zerolog.Logger
}

// getZerologLevel parses a string level into a zerolog.Level.
func getZerologLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info", "":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// NewZerologLogger creates a zerolog logger.
// env: "development" -> pretty console; else JSON.
func NewZerologLogger(env string) Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var z zerolog.Logger
	var lvl zerolog.Level // Instance-level

	if env == "development" {
		cw := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339Nano,
		}
		z = zerolog.New(cw).
			With().
			Timestamp().
			Caller().
			Logger()
		lvl = zerolog.DebugLevel // Default for dev
	} else {
		z = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Logger()
		lvl = zerolog.InfoLevel // Default for prod
	}

	// Set the level on the INSTANCE, not globally
	z = z.Level(lvl)

	return &ZerologAdapter{logger: z}
}

// SetLevel updates the log level for this specific logger instance.
func (l *ZerologAdapter) SetLevel(level string) {
	parsedLevel := getZerologLevel(level)
	// Create a new logger with the new level and assign it
	l.logger = l.logger.Level(parsedLevel)
}

func (l *ZerologAdapter) With(f Fields) Logger {
	ctx := l.logger.With()
	for k, v := range f {
		ctx = ctx.Interface(k, v)
	}
	nl := ctx.Logger()
	return &ZerologAdapter{logger: nl}
}

// Plain methods
func (l *ZerologAdapter) Info(msg string) { l.logger.Info().Msg(msg) }

// Error implements Logger.
func (l *ZerologAdapter) Error(err error, msg string) {
	l.logger.Error().Err(err).Msg(msg)
}
func (l *ZerologAdapter) Warn(msg string)  { l.logger.Warn().Msg(msg) }
func (l *ZerologAdapter) Debug(msg string) { l.logger.Debug().Msg(msg) }
func (l *ZerologAdapter) Fatal(msg string) { l.logger.Fatal().Msg(msg) }

package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Fields map[string]any

type Logger interface {
	Info(msg string)
	Error(err error, msg string)
	Warn(msg string)
	Debug(msg string)
	Fatal(msg string)

	// Tambahan untuk structured logging & konfigurasi:
	With(f Fields) Logger
	SetLevel(level string) // "debug", "info", "warn", "error"
}

// ZerologAdapter implements Logger.
type ZerologAdapter struct {
	logger zerolog.Logger
}

// NewZerologLogger creates a zerolog logger.
// env: "development" -> pretty console; else JSON.
func NewZerologLogger(env string) Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var z zerolog.Logger
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
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		z = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Logger()
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return &ZerologAdapter{logger: z}
}

func (l *ZerologAdapter) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info", "":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func (l *ZerologAdapter) With(f Fields) Logger {
	ctx := l.logger.With()
	for k, v := range f {
		ctx = ctx.Interface(k, v)
	}
	nl := ctx.Logger()
	return &ZerologAdapter{logger: nl}
}

// Plain methods (backward-compatible)
func (l *ZerologAdapter) Info(msg string) { l.logger.Info().Msg(msg) }
func (l *ZerologAdapter) Error(err error, msg string) {
	l.logger.Error().Err(err).Msg(msg)
}
func (l *ZerologAdapter) Warn(msg string)  { l.logger.Warn().Msg(msg) }
func (l *ZerologAdapter) Debug(msg string) { l.logger.Debug().Msg(msg) }
func (l *ZerologAdapter) Fatal(msg string) { l.logger.Fatal().Msg(msg) }

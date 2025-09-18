package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(msg string)
	Error(err error, msg string)
	Warn(msg string)
	Debug(msg string)
	Fatal(msg string)
}

// ZerologAdapter implements the ports.Logger interface.
type ZerologAdapter struct {
	logger zerolog.Logger
}

// NewZerologLogger creates a new zerolog logger instance.
func NewZerologLogger(env string) Logger {
	// Configure zerolog to output to stdout with color-coded, human-readable format
	// for development and plain JSON for production.
	var zlog zerolog.Logger
	if env == "development" {
		zlog = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	} else {
		zlog = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Set the global log level based on the environment.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return &ZerologAdapter{logger: zlog}
}

// Below are the implementations of the ports.Logger methods.
func (l *ZerologAdapter) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *ZerologAdapter) Error(err error, msg string) {
	l.logger.Error().Err(err).Msg(msg)
}

func (l *ZerologAdapter) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l *ZerologAdapter) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l *ZerologAdapter) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

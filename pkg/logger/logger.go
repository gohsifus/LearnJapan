package logger

import (
	"LearnJapan.com/configs"
	"LearnJapan.com/constants"
	"github.com/rs/zerolog"
	"os"
)

type Logger struct {
	cfg       *configs.Configs
	stdLogger zerolog.Logger
}

func NewLogger(cfg *configs.Configs) *Logger {
	return &Logger{
		cfg:       cfg,
		stdLogger: zerolog.New(os.Stdout),
	}
}

func (l Logger) Error(message interface{}) {
	switch log := message.(type) {
	case error:
		l.stdLogger.Error().Msg(log.Error())
	case string:
		l.stdLogger.Error().Msg(log)
	default:
		l.stdLogger.Error().Msg(constants.MESSAGE_UNKNOWN_LOG_TYPE)
	}
}

func (l Logger) Info(message string) {
	l.stdLogger.Info().Msg(message)
}

func (l Logger) Fatal(message interface{}) {
	switch log := message.(type) {
	case error:
		l.stdLogger.Fatal().Msg(log.Error())
	case string:
		l.stdLogger.Fatal().Msg(log)
	default:
		l.stdLogger.Fatal().Msg(constants.MESSAGE_UNKNOWN_LOG_TYPE)
	}
}

func (l Logger) Warn(message string) {
	l.stdLogger.Warn().Msg(message)
}

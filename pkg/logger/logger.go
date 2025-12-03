package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	// Настройка логгера
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	// Установка уровня логирования из переменных окружения
	if level, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		zerolog.SetGlobalLevel(level)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return logger
}

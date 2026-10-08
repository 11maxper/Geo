package logger

import (
	"log"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	return logger
}

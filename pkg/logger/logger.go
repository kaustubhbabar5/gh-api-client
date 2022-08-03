package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// builds and returns logger
func New() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction :%w", err)
	}

	return logger, nil
}

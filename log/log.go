package log

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func Logger() *zap.Logger {
	return logger
}

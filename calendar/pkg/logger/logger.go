package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Create creates new Zap logger
func Create(filePath string, level string) (*zap.Logger, error) {
	var zapLevel zapcore.Level
	err := zapLevel.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zapLevel)
	loggerConfig.OutputPaths = []string{filePath}
	loggerConfig.ErrorOutputPaths = []string{filePath}

	return loggerConfig.Build()
}

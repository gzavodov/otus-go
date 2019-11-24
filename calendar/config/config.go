package config

import (
	"go.uber.org/zap/zapcore"
)

//Configuration Web server configuration
type Configuration struct {
	HTTPAddress           string `json:"http_listen"`
	LogFilePath           string `json:"log_file"`
	LogLevel              string `json:"log_level"`
	EventRepositoryTypeID int    `json:"event_repository_type_id"`
}

//ParseZapLogLevel Parse Zap LogLevel from strig
func (c *Configuration) ParseZapLogLevel() (zapcore.Level, error) {
	var l zapcore.Level
	err := l.UnmarshalText([]byte(c.LogLevel))
	return l, err
}

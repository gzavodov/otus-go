package web

import (
	"go.uber.org/zap/zapcore"
)

//Config Web server configuration
type Config struct {
	HTTPAddress string `yaml:"http_listen"`
	LogFilePath string `yaml:"log_file"`
	LogLevel    string `yaml:"log_level"`
}

//ParseZapLogLevel Parse Zap LogLevel from strig
func (c *Config) ParseZapLogLevel() (zapcore.Level, error) {
	var l zapcore.Level
	err := l.UnmarshalText([]byte(c.LogLevel))
	return l, err
}
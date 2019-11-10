package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/gzavodov/otus-go/calendar/app/web"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTPAddress string `yaml:"http_listen"`
	LogFilePath string `yaml:"log_file"`
	LogLevel    string `yaml:"log_level"`
}

func (c *Config) ParseZapLogLevel() (zapcore.Level, error) {
	var l zapcore.Level
	err := l.UnmarshalText([]byte(c.LogLevel))
	return l, err
}

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "config.yml"
	}

	configFile, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatalf("Could not read configuration file: %v", err)
	}

	var config Config
	if yaml.Unmarshal(configFile, &config) != nil {
		log.Fatalf("Could not internalize configuration file data: %v", err)
	}

	if config.HTTPAddress == "" {
		config.HTTPAddress = "127.0.0.1:8080"
	}

	if config.LogFilePath == "" {
		config.LogFilePath = "stderr"
	}

	if config.LogLevel == "" {
		config.LogLevel = "debug"
	}

	logLelev, err := config.ParseZapLogLevel()
	if err != nil {
		log.Fatalf("Could not internalize zap logger level: %v", err)
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(logLelev)
	loggerConfig.OutputPaths = []string{config.LogFilePath}
	loggerConfig.ErrorOutputPaths = []string{config.LogFilePath}

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("Could not initialize zap logger: %v", err)
	}
	defer logger.Sync()

	server := web.Server{Logger: logger}
	err = server.Start(config.HTTPAddress)
	if err != nil {
		logger.Fatal("Could not start HTTP server", zap.Error(err))
	}
}

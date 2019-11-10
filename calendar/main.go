package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/gzavodov/otus-go/calendar/app/web"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "http_config.yml"
	}

	configFile, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatalf("Could not read configuration file: %v", err)
	}

	config := &web.Config{}
	if yaml.Unmarshal(configFile, config) != nil {
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

	server := web.NewServer(config, logger)
	err = server.Start(config.HTTPAddress)
	if err != nil {
		logger.Fatal("Could not start HTTP server", zap.Error(err))
	}
}

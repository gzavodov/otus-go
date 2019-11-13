package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/gzavodov/otus-go/calendar/app/web"
	"github.com/gzavodov/otus-go/calendar/config"
	"go.uber.org/zap"
)

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "./config/config.development.json"
	}

	configFile, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatalf("Could not read configuration file: %v", err)
	}

	config := &config.Configuration{}
	if json.Unmarshal(configFile, config) != nil {
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

	log.Printf("Starting web server on %s\n", config.HTTPAddress)

	server := web.NewServer(config.HTTPAddress, logger)
	err = server.Start()
	if err != nil {
		log.Fatalf("Could not initialize zap logger: %v", err)
	}
}

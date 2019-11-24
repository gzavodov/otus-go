package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
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

	configuration := &config.Configuration{}
	if json.Unmarshal(configFile, configuration) != nil {
		log.Fatalf("Could not internalize configuration file data: %v", err)
	}

	if configuration.HTTPAddress == "" {
		configuration.HTTPAddress = "127.0.0.1:8080"
	}

	if configuration.LogFilePath == "" {
		configuration.LogFilePath = "stderr"
	}

	if configuration.LogLevel == "" {
		configuration.LogLevel = "debug"
	}

	logLelev, err := configuration.ParseZapLogLevel()
	if err != nil {
		log.Fatalf("Could not internalize zap logger level: %v", err)
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(logLelev)
	loggerConfig.OutputPaths = []string{configuration.LogFilePath}
	loggerConfig.ErrorOutputPaths = []string{configuration.LogFilePath}

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("Could not initialize zap logger: %v", err)
	}
	defer logger.Sync()

	log.Printf("Starting web server on %s\n", configuration.HTTPAddress)

	repo, err := repository.CreateEventRepository(configuration.EventRepositoryTypeID)
	if err != nil {
		log.Fatalf("Could not create event repository: %v", err)
	}
	server := web.NewServer(configuration.HTTPAddress, repo, logger)
	err = server.Start()
	if err != nil {
		log.Fatalf("Could not start web server: %v", err)
	}
}

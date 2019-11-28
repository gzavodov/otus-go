package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/gzavodov/otus-go/calendar/app/factory"
	"github.com/gzavodov/otus-go/calendar/app/logger"
	"github.com/gzavodov/otus-go/calendar/config"
)

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "./config/config.development.rpc.json"
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
		configuration.HTTPAddress = "127.0.0.1:9090"
	}

	if configuration.LogFilePath == "" {
		configuration.LogFilePath = "stderr"
	}

	if configuration.LogLevel == "" {
		configuration.LogLevel = "debug"
	}

	appLogger, err := logger.Create(configuration.LogFilePath, configuration.LogLevel)
	if err != nil {
		log.Fatalf("Could not initialize zap logger: %v", err)
	}
	defer appLogger.Sync()

	appRepo, err := factory.CreateEventRepository(configuration.EventRepositoryTypeID)
	if err != nil {
		log.Fatalf("Could not create event repository: %v", err)
	}

	service, err := factory.CreateEndpointService(
		configuration.EndpointServiceTypeID,
		configuration.HTTPAddress,
		appRepo,
		appLogger,
	)
	if err != nil {
		log.Fatalf("Could not create endpoint service: %v", err)
	}

	log.Printf("Starting %s service on %s...\n", service.GetServiceName(), configuration.HTTPAddress)

	err = service.Start()
	if err != nil {
		log.Fatalf("Could not start endpoint service: %v", err)
	}
}

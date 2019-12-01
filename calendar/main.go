package main

import (
	"context"
	"flag"
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

	configuration := &config.Configuration{}
	err := configuration.Load(
		*configFilePath,
		&config.Configuration{
			HTTPAddress: "127.0.0.1:9090",
			LogFilePath: "stderr",
			LogLevel:    "debug",
		},
	)
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	appLogger, err := logger.Create(configuration.LogFilePath, configuration.LogLevel)
	if err != nil {
		log.Fatalf("Could not initialize zap logger: %v", err)
	}
	defer appLogger.Sync()

	appRepo, err := factory.CreateEventRepository(
		context.Background(),
		configuration.EventRepositoryTypeID,
		configuration.EventRepositoryDSN,
	)
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

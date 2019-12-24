package main

import (
	"context"
	"flag"
	"log"

	"github.com/gzavodov/otus-go/calendar/factory/repofactory"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
	"github.com/gzavodov/otus-go/calendar/pkg/logger"
	"github.com/gzavodov/otus-go/calendar/service/rpc"
)

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "./config/config.development.json"
	}

	configuration := &config.Configuration{}
	err := configuration.Load(
		*configFilePath,
		&config.Configuration{
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

	appRepo, err := repofactory.CreateEventRepository(
		context.Background(),
		configuration.EventRepositoryTypeID,
		configuration.EventRepositoryDSN,
	)

	if err != nil {
		log.Fatalf("Could not create event repository: %v", err)
	}

	appService := rpc.NewServer(configuration.HTTPAddress, appRepo, appLogger)

	log.Printf("Starting %s service on %s...\n", appService.GetServiceName(), configuration.HTTPAddress)
	if err = appService.Start(); err != nil {
		log.Fatalf("Could not start endpoint service: %v", err)
	}
}

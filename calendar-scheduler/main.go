package main

import (
	"context"
	"flag"
	"log"

	"github.com/gzavodov/otus-go/calendar/factory/queuefactory"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
	"github.com/gzavodov/otus-go/calendar/pkg/logger"
	"github.com/gzavodov/otus-go/calendar/service/rpc"
	"github.com/gzavodov/otus-go/calendar/service/scheduler"
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

	queueChannel, err := queuefactory.CreateQueueChannel(
		context.Background(),
		configuration.AMPQTypeID,
		configuration.AMPQName,
		configuration.AMPQAddress,
	)

	if err != nil {
		log.Fatalf("Could not create queue channel: %v", err)
	}

	appService := scheduler.NewServer(
		context.Background(),
		queueChannel,
		rpc.NewEventRepository(context.Background(), configuration.HTTPAddress),
		configuration.SchedulerCheckInterval,
		appLogger,
	)

	log.Printf("Starting sheduler server on queue %s on %s...\n", configuration.AMPQName, configuration.AMPQAddress)
	if err = appService.Start(); err != nil {
		log.Fatalf("Could not start scheduler server: %v", err)
	}
}

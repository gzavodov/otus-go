package main

import (
	"context"
	"flag"
	"log"
	"otus-go/calendar-scheduler/config"

	//_ "github.com/gzavodov/otus-go/calendar-scheduler/config"
	"github.com/gzavodov/otus-go/calendar/app/factory"
	"github.com/gzavodov/otus-go/calendar/app/logger"
	"github.com/gzavodov/otus-go/calendar/app/queuefactory"
	"github.com/gzavodov/otus-go/calendar/app/scheduler"
)

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "./config/config.development.ampq.json"
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

	appRepo, err := factory.CreateEventRepository(
		context.Background(),
		configuration.EventRepositoryTypeID,
		configuration.EventRepositoryDSN,
	)
	if err != nil {
		log.Fatalf("Could not create event repository: %v", err)
	}

	queueChannel, err := queuefactory.CreateQueueChannel(
		context.Background(),
		configuration.AMPQTypeID,
		configuration.AMPQName,
		configuration.AMPQAddress,
	)

	if err != nil {
		log.Fatalf("Could not create queue channel: %v", err)
	}

	server := scheduler.NewServer(
		context.Background(),
		queueChannel,
		appRepo,
		configuration.SchedulerCheckInterval,
		appLogger,
	)

	log.Printf("Starting sheduler server on queue %s on %s...\n", configuration.AMPQName, configuration.AMPQAddress)

	err = server.Start()
	if err != nil {
		log.Fatalf("Could not start scheduler server: %v", err)
	}
}

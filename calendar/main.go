package main

import (
	"context"
	"flag"
	"log"

	"github.com/gzavodov/otus-go/calendar/app/factory"
	"github.com/gzavodov/otus-go/calendar/app/logger"
	"github.com/gzavodov/otus-go/calendar/app/queuefactory"
	"github.com/gzavodov/otus-go/calendar/app/scheduler"
	"github.com/gzavodov/otus-go/calendar/config"
)

const (
	modeWeb        = "web"
	modeRPC        = "rpc"
	modeAMPQWriter = "ampq_writer"
	modeAMPQReader = "ampq_reader"
)

func main() {
	mode := flag.String("mode", "", "Application mode (web, rpc, ampq_writer, ampq_reader)")
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *mode == "" {
		*mode = modeRPC
	}

	if *configFilePath == "" {
		if *mode == modeRPC {
			*configFilePath = "./config/config.development.rpc.json"
		} else if *mode == modeWeb {
			*configFilePath = "./config/config.development.web.json"
		} else if *mode == modeAMPQWriter || *mode == modeAMPQReader {
			*configFilePath = "./config/config.development.ampq.json"
		} else {
			log.Fatal("Could not find configuration file path")
		}
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

	if *mode == modeWeb || *mode == modeRPC {
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
	} else if *mode == modeAMPQWriter {
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
	} else if *mode == modeAMPQReader {
		queueChannel, err := queuefactory.CreateQueueChannel(
			context.Background(),
			configuration.AMPQTypeID,
			configuration.AMPQName,
			configuration.AMPQAddress,
		)

		if err != nil {
			log.Fatalf("Could not create queue channel: %v", err)
		}

		client := scheduler.NewClient(
			context.Background(),
			queueChannel,
			appLogger,
		)

		log.Printf("Starting sheduler client on queue %s on %s...\n", configuration.AMPQName, configuration.AMPQAddress)

		err = client.Start()
		if err != nil {
			log.Fatalf("Could not start scheduler client: %v", err)
		}
	}
}

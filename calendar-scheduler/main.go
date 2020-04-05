package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/gzavodov/otus-go/calendar/factory/queuefactory"
	"github.com/gzavodov/otus-go/calendar/factory/repofactory"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
	"github.com/gzavodov/otus-go/calendar/pkg/logger"
	"github.com/gzavodov/otus-go/calendar/pkg/queuemonitoring"
	"github.com/gzavodov/otus-go/calendar/service/scheduler"
	"github.com/gzavodov/otus-go/calendar/service/sysmonitor"
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

	defer func() {
		if err := appLogger.Sync(); err != nil {
			log.Fatalf("could not flush logger write buffers: %v", err)
		}
	}()

	appRepo, err := repofactory.CreateEventRepository(
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

	queueMonitoring := queuemonitoring.NewMiddleware("api", appLogger)

	wg := &sync.WaitGroup{}

	if configuration.HealthcheckHTTPAddress != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sysMonitoringService := sysmonitor.NewServer(configuration.HealthcheckHTTPAddress, queueMonitoring, appLogger)
			log.Printf("Starting %s service on %s...\n", sysMonitoringService.GetServiceName(), configuration.HealthcheckHTTPAddress)

			if err = sysMonitoringService.Start(); err != nil {
				log.Fatalf("Could not start System Monitoring Service: %v", err)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		appService := scheduler.NewServer(
			context.Background(),
			queueChannel,
			appRepo,
			configuration.SchedulerCheckInterval,
			appLogger,
		)
		appService.RegisterMonitoringMiddleware(queueMonitoring)

		log.Printf("Starting sheduler server on queue %s on %s...\n", configuration.AMPQName, configuration.AMPQAddress)
		if err = appService.Start(); err != nil {
			log.Fatalf("Could not start scheduler server: %v", err)
		}
	}()

	wg.Wait()
}

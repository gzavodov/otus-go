package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/gzavodov/otus-go/calendar/factory/repofactory"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
	"github.com/gzavodov/otus-go/calendar/pkg/httpmonitoring"
	"github.com/gzavodov/otus-go/calendar/pkg/logger"
	"github.com/gzavodov/otus-go/calendar/service/sysmonitor"
	"github.com/gzavodov/otus-go/calendar/service/web"
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

	webMonitoring := httpmonitoring.NewMiddleware("api", appLogger)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		sysMonitoringService := sysmonitor.NewServer(configuration.HealthcheckHTTPAddress, webMonitoring, appLogger)
		log.Printf("Starting %s service on %s...\n", sysMonitoringService.GetServiceName(), configuration.HealthcheckHTTPAddress)

		if err = sysMonitoringService.Start(); err != nil {
			log.Fatalf("Could not start System Monitoring Service: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		appService := web.NewServer(configuration.HTTPAddress, appRepo, appLogger)
		appService.MonitoringMiddleware = webMonitoring
		log.Printf("Starting %s service on %s...\n", appService.GetServiceName(), configuration.HTTPAddress)
		if err = appService.Start(); err != nil {
			log.Fatalf("Could not start RPC Service: %v", err)
		}
	}()

	wg.Wait()
}

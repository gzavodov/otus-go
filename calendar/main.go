package main

import (
	"context"
	"flag"
	"log"

	"github.com/gzavodov/otus-go/calendar/pkg/config"
	"github.com/gzavodov/otus-go/calendar/pkg/logger"
	"github.com/gzavodov/otus-go/calendar/service"
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

	defer func() {
		if err := appLogger.Sync(); err != nil {
			log.Fatalf("could not flush logger write buffers: %v", err)
		}
	}()

	if *mode == modeWeb || *mode == modeRPC || *mode == modeAMPQWriter {
		appService, err := service.CreateService(context.Background(), configuration, appLogger)
		if err != nil {
			log.Fatalf("Could not create service: %v", err)
		}

		log.Printf("Starting %s service...\n", appService.GetServiceName())
		if err := appService.Start(); err != nil {
			log.Fatalf("Could not start service: %v", err)
		}
	}
}

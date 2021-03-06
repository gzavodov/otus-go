package main

import (
	"context"
	"flag"
	"log"

	"github.com/gzavodov/otus-go/banner-rotation/algorithm"
	"github.com/gzavodov/otus-go/banner-rotation/config"
	"github.com/gzavodov/otus-go/banner-rotation/internal/rabbitmq"
	"github.com/gzavodov/otus-go/banner-rotation/internal/sql"
	"github.com/gzavodov/otus-go/banner-rotation/logger"
	"github.com/gzavodov/otus-go/banner-rotation/rest"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "./config/config.development.json"
	}

	conf := &config.Configuration{}
	err := conf.Load(
		*configFilePath,
		&config.Configuration{
			AlgorithmTypeID: algorithm.TypeUCB1,
			LogFilePath:     "stderr",
			LogLevel:        "debug",
		},
	)
	if err != nil {
		log.Fatalf("could not load configuration: %v", err)
	}

	appLogger, err := logger.Create(conf.LogFilePath, conf.LogLevel)
	if err != nil {
		log.Fatalf("could not initialize logger: %v", err)
	}
	defer func() {
		if err := appLogger.Sync(); err != nil {
			log.Fatalf("could not flush logger write buffers: %v", err)
		}
	}()

	appContext := context.Background()

	bannerRepo := sql.NewBannerRepository(appContext, conf.RepositoryDSN)
	slotRepo := sql.NewSlotRepository(appContext, conf.RepositoryDSN)
	bindingRepo := sql.NewBindingRepository(appContext, conf.RepositoryDSN)
	groupRepo := sql.NewGroupRepository(appContext, conf.RepositoryDSN)
	statisticsRepo := sql.NewStatisticsRepository(appContext, conf.RepositoryDSN)

	appService := rest.NewServer(
		conf.HTTPAddress,
		usecase.NewBannerUsecase(bannerRepo, bindingRepo, statisticsRepo, conf.AlgorithmTypeID),
		usecase.NewSlotUsecase(slotRepo),
		usecase.NewGroupUsecase(groupRepo),
		rabbitmq.NewChannel(appContext, conf.AMPQName, conf.AMPQAddress),
		appLogger,
	)

	log.Printf("Starting %s service...\n", appService.GetServiceName())
	if err := appService.Start(); err != nil {
		log.Fatalf("could not start %s service: %v", appService.GetServiceName(), err)
	}
}

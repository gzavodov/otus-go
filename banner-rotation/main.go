package main

import (
	"context"
	"flag"
	"log"

	"github.com/gzavodov/otus-go/banner-rotation/config"
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
			LogFilePath: "stderr",
			LogLevel:    "debug",
		},
	)
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	appLogger, err := logger.Create(conf.LogFilePath, conf.LogLevel)
	if err != nil {
		log.Fatalf("Could not initialize zap logger: %v", err)
	}
	defer appLogger.Sync()

	appContext := context.Background()

	bannerRepo := sql.NewBannerRepository(appContext, conf.RepositoryDSN)
	slotRepo := sql.NewSlotRepository(appContext, conf.RepositoryDSN)
	bindingRepo := sql.NewBindingRepository(appContext, conf.RepositoryDSN)
	groupRepo := sql.NewGroupRepository(appContext, conf.RepositoryDSN)
	statisticsRepo := sql.NewStatisticsRepository(appContext, conf.RepositoryDSN)

	appService := rest.NewServer(
		conf.HTTPAddress,
		usecase.NewBannerUsecase(bannerRepo, bindingRepo, statisticsRepo),
		usecase.NewSlotUsecase(slotRepo),
		usecase.NewGroupUsecase(groupRepo),
		appLogger,
	)

	log.Printf("Starting %s service...\n", appService.GetServiceName())
	if err := appService.Start(); err != nil {
		log.Fatalf("Could not start %s service: %v", appService.GetServiceName(), err)
	}
}

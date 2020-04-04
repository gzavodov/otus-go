package service

import (
	"context"
	"fmt"

	"github.com/gzavodov/otus-go/calendar/factory/queuefactory"
	"github.com/gzavodov/otus-go/calendar/factory/repofactory"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
	"github.com/gzavodov/otus-go/calendar/pkg/service"
	"github.com/gzavodov/otus-go/calendar/service/rpc"
	"github.com/gzavodov/otus-go/calendar/service/scheduler"
	"github.com/gzavodov/otus-go/calendar/service/web"
	"go.uber.org/zap"
)

//Service Type
const (
	TypeUnknown   = 0
	TypeWeb       = 1
	TypeGRPC      = 2
	TypeScheduler = 3
	TypeClient    = 4
)

//CreateService creates service by configuration
func CreateService(ctx context.Context, conf *config.Configuration, logger *zap.Logger) (service.Service, error) {
	switch conf.ServiceTypeID {
	case TypeWeb:
		repo, err := repofactory.CreateEventRepository(
			ctx,
			conf.EventRepositoryTypeID,
			conf.EventRepositoryDSN,
		)

		if err != nil {
			return nil, fmt.Errorf("could not create repository (%w)", err)
		}

		return web.NewServer(conf.GRPCAddress, repo, logger), nil
	case TypeGRPC:
		repo, err := repofactory.CreateEventRepository(
			ctx,
			conf.EventRepositoryTypeID,
			conf.EventRepositoryDSN,
		)

		if err != nil {
			return nil, fmt.Errorf("could not create repository (%w)", err)
		}

		return rpc.NewServer(conf.HTTPAddress, repo, logger), nil
	case TypeScheduler:
		repo, err := repofactory.CreateEventRepository(
			ctx,
			conf.EventRepositoryTypeID,
			conf.EventRepositoryDSN,
		)

		if err != nil {
			return nil, fmt.Errorf("could not create repository (%w)", err)
		}

		queueChannel, err := queuefactory.CreateQueueChannel(
			ctx,
			conf.AMPQTypeID,
			conf.AMPQName,
			conf.AMPQAddress,
		)

		if err != nil {
			return nil, fmt.Errorf("could not create queue channel (%w)", err)
		}

		return scheduler.NewServer(
				ctx,
				queueChannel,
				repo,
				conf.SchedulerCheckInterval,
				logger,
			),
			nil
	case TypeClient:
		queueChannel, err := queuefactory.CreateQueueChannel(
			ctx,
			conf.AMPQTypeID,
			conf.AMPQName,
			conf.AMPQAddress,
		)

		if err != nil {
			return nil, fmt.Errorf("could not create queue channel (%w)", err)
		}

		return scheduler.NewClient(
				ctx,
				queueChannel,
				nil,
				logger,
			),
			nil

	default:
		return nil, fmt.Errorf("service type %d is not supported in current context", conf.ServiceTypeID)
	}
}

package factory

import (
	"fmt"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"github.com/gzavodov/otus-go/calendar/app/endpoint"
	"github.com/gzavodov/otus-go/calendar/app/rpc"
	"github.com/gzavodov/otus-go/calendar/app/web"
	"go.uber.org/zap"
)

//Endpoint service Type
const (
	EndpointTypeUnknown = 0
	EndpointTypeWeb     = 1
	EndpointTypeGRPC    = 2
)

//CreateEndpointService creates endpoint service by type
func CreateEndpointService(typeID int, address string, repo repository.EventRepository, logger *zap.Logger) (endpoint.Service, error) {
	switch typeID {
	case EndpointTypeWeb:
		return web.NewServer(address, repo, logger), nil
	case EndpointTypeGRPC:
		return rpc.NewServer(address, repo, logger), nil
	default:
		return nil, fmt.Errorf("endpoint service type %d is not supported in current context", typeID)
	}
}

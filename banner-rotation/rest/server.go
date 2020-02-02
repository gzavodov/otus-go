package rest

import (
	"context"
	"net/http"

	"github.com/gzavodov/otus-go/banner-rotation/endpoint"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
	"go.uber.org/zap"
)

//NewServer Creates new Web server
func NewServer(
	address string,
	bannerUsecase *usecase.Banner,
	slotUsecase *usecase.Slot,
	bindingUsecase *usecase.Binding,
	groupUsecase *usecase.Group,
	logger *zap.Logger,
) *Server {
	return &Server{
		Server:         endpoint.Server{Name: "REST", Address: address, Logger: logger},
		bannerUsecase:  bannerUsecase,
		slotUsecase:    slotUsecase,
		bindingUsecase: bindingUsecase,
		groupUsecase:   groupUsecase,
	}
}

//Server Simple Web Server for calendar event API
type Server struct {
	server *http.Server

	bannerHandler *Banner
	bannerUsecase *usecase.Banner

	slotHandler *Slot
	slotUsecase *usecase.Slot

	bindingHandler *Binding
	bindingUsecase *usecase.Binding

	groupHandler *Group
	groupUsecase *usecase.Group

	endpoint.Server
}

//Start Start handling of Web requests
func (s *Server) Start() error {
	serverMux := http.NewServeMux()

	s.bannerHandler = &Banner{
		Handler: endpoint.Handler{
			Name:        "Banner",
			ServiceName: s.Name,
			Logger:      s.Logger,
		},
		ucase: s.bannerUsecase,
	}
	serverMux.HandleFunc("/banner/create", s.bannerHandler.Create)
	serverMux.HandleFunc("/banner/read", s.bannerHandler.Read)
	serverMux.HandleFunc("/banner/update", s.bannerHandler.Update)
	serverMux.HandleFunc("/banner/delete", s.bannerHandler.Delete)

	serverMux.HandleFunc("/banner/add-to-slot", s.bannerHandler.AddToSlot)
	serverMux.HandleFunc("/banner/delete-from-slot", s.bannerHandler.DeleteFromSlot)
	serverMux.HandleFunc("/banner/register-click", s.bannerHandler.RegisterClick)
	serverMux.HandleFunc("/banner/choose", s.bannerHandler.Choose)

	s.slotHandler = &Slot{
		Handler: endpoint.Handler{
			Name:        "Slot",
			ServiceName: s.Name,
			Logger:      s.Logger,
		},
		ucase: s.slotUsecase,
	}
	serverMux.HandleFunc("/slot/create", s.slotHandler.Create)
	serverMux.HandleFunc("/slot/read", s.slotHandler.Read)
	serverMux.HandleFunc("/slot/update", s.slotHandler.Update)
	serverMux.HandleFunc("/slot/delete", s.slotHandler.Delete)

	s.bindingHandler = &Binding{
		Handler: endpoint.Handler{
			Name:        "Binding",
			ServiceName: s.Name,
			Logger:      s.Logger,
		},
		ucase: s.bindingUsecase,
	}
	serverMux.HandleFunc("/binding/create", s.bindingHandler.Create)
	serverMux.HandleFunc("/binding/read", s.bindingHandler.Read)
	serverMux.HandleFunc("/binding/update", s.bindingHandler.Update)
	serverMux.HandleFunc("/binding/delete", s.bindingHandler.Delete)

	s.groupHandler = &Group{
		Handler: endpoint.Handler{
			Name:        "Group",
			ServiceName: s.Name,
			Logger:      s.Logger,
		},
		ucase: s.groupUsecase,
	}
	serverMux.HandleFunc("/group/create", s.groupHandler.Create)
	serverMux.HandleFunc("/group/read", s.groupHandler.Read)
	serverMux.HandleFunc("/group/update", s.groupHandler.Update)
	serverMux.HandleFunc("/group/delete", s.groupHandler.Delete)

	s.server = &http.Server{Addr: s.Address, Handler: serverMux}

	err := s.server.ListenAndServe()
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return err
}

//Stop stop handling of Web requests
func (s *Server) Stop() {
	if s.server != nil {
		s.server.Shutdown(context.Background())
	}
}

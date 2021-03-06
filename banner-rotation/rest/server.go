package rest

import (
	"context"
	"net/http"

	//enabling pprof debugger
	"net/http/pprof"
	"github.com/gzavodov/otus-go/banner-rotation/queue"

	"github.com/gzavodov/otus-go/banner-rotation/endpoint"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
	"go.uber.org/zap"
)

const PostMethod = "POST"

//NewServer Creates new Web server
func NewServer(
	address string,

	bannerUsecase *usecase.Banner,
	slotUsecase *usecase.Slot,
	groupUsecase *usecase.Group,

	notificationChannel queue.NotificationChannel,

	logger *zap.Logger,
) *Server {
	return &Server{
		Server:              endpoint.Server{Name: "REST", Address: address, Logger: logger},
		bannerUsecase:       bannerUsecase,
		slotUsecase:         slotUsecase,
		groupUsecase:        groupUsecase,
		notificationChannel: notificationChannel,
	}
}

//Server Simple Web Server for calendar event API
type Server struct {
	server *http.Server

	bannerHandler *BannerHandler
	bannerUsecase *usecase.Banner

	slotHandler *EntityHandler
	slotUsecase *usecase.Slot

	groupHandler *EntityHandler
	groupUsecase *usecase.Group

	notificationChannel queue.NotificationChannel

	endpoint.Server
}

//Start Start handling of Web requests
func (s *Server) Start() error {
	serverMux := http.NewServeMux()

	s.bannerHandler = NewBannerHandler(s.bannerUsecase, s.Name, s.notificationChannel, s.Logger)

	serverMux.HandleFunc("/banner/create", s.bannerHandler.Create)
	serverMux.HandleFunc("/banner/read", s.bannerHandler.Read)
	serverMux.HandleFunc("/banner/update", s.bannerHandler.Update)
	serverMux.HandleFunc("/banner/delete", s.bannerHandler.Delete)
	serverMux.HandleFunc("/banner/get-by-caption", s.bannerHandler.GetByCaption)

	serverMux.HandleFunc("/banner/add-to-slot", s.bannerHandler.AddToSlot)
	serverMux.HandleFunc("/banner/delete-from-slot", s.bannerHandler.DeleteFromSlot)
	serverMux.HandleFunc("/banner/is-in-slot", s.bannerHandler.IsInSlot)

	serverMux.HandleFunc("/banner/register-click", s.bannerHandler.RegisterClick)
	serverMux.HandleFunc("/banner/choose", s.bannerHandler.Choose)

	s.slotHandler = &EntityHandler{
		Accessor: &Slot{ucase: s.slotUsecase},
		Handler: endpoint.Handler{
			Name:                "Slot",
			ServiceName:         s.Name,
			Logger:              s.Logger,
			NotificationChannel: s.notificationChannel,
		},
	}
	serverMux.HandleFunc("/slot/create", s.slotHandler.Create)
	serverMux.HandleFunc("/slot/read", s.slotHandler.Read)
	serverMux.HandleFunc("/slot/update", s.slotHandler.Update)
	serverMux.HandleFunc("/slot/delete", s.slotHandler.Delete)
	serverMux.HandleFunc("/slot/get-by-caption", s.slotHandler.GetByCaption)

	s.groupHandler = &EntityHandler{
		Accessor: &Group{ucase: s.groupUsecase},
		Handler: endpoint.Handler{
			Name:                "Group",
			ServiceName:         s.Name,
			Logger:              s.Logger,
			NotificationChannel: s.notificationChannel,
		},
	}
	serverMux.HandleFunc("/group/create", s.groupHandler.Create)
	serverMux.HandleFunc("/group/read", s.groupHandler.Read)
	serverMux.HandleFunc("/group/update", s.groupHandler.Update)
	serverMux.HandleFunc("/group/delete", s.groupHandler.Delete)
	serverMux.HandleFunc("/group/get-by-caption", s.groupHandler.GetByCaption)

	serverMux.HandleFunc("/debug/pprof/", pprof.Index)
	serverMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	serverMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	serverMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	s.server = &http.Server{Addr: s.Address, Handler: serverMux}

	err := s.server.ListenAndServe()
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return err
}

//Stop stop handling of Web requests
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}

	return nil
}

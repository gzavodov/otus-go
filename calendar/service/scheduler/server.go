package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/gzavodov/otus-go/calendar/model"
	"github.com/gzavodov/otus-go/calendar/pkg/queue"
	"github.com/gzavodov/otus-go/calendar/repository"
	"go.uber.org/zap"
)

//NewServer creates new scheduler server
func NewServer(ctx context.Context, channel queue.NotificationChannel, repo repository.EventRepository, checkInterval int, logger *zap.Logger) *Server {
	if checkInterval <= 0 {
		checkInterval = 60
	}
	return &Server{
		ctx:             ctx,
		channel:         channel,
		shutdownChannel: make(chan struct{}),
		repo:            repo,
		checkInterval:   checkInterval,
		logger:          logger,
	}
}

//Server event reminder service
type Server struct {
	ctx             context.Context
	channel         queue.NotificationChannel
	shutdownChannel chan struct{}
	repo            repository.EventRepository
	once            sync.Once
	checkInterval   int
	logger          *zap.Logger
}

//GetServiceName returns service name
func (s *Server) GetServiceName() string {
	return "Scheduler server"
}

//Start starts scheduler server
func (s *Server) Start() error {
	err := s.Check()
	if err != nil {
		s.LogError(ErrorCategoryRepository, err)
		return err
	}

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-s.shutdownChannel:
			return nil
		case <-time.After(time.Duration(s.checkInterval) * time.Second):
			err = s.Check()
			if err != nil {
				s.LogError(ErrorCategoryRepository, err)
				return err
			}
		}
	}
}

//Stop stop scheduler server
func (s *Server) Stop() {
	s.once.Do(func() { s.shutdownChannel <- struct{}{} })
}

//Check checks events and sends notification if required
func (s *Server) Check() error {
	list, err := s.repo.ReadNotificationList(0, time.Now())
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return nil
	}

	notifications := make([]*queue.Notification, 0, len(list))
	for _, item := range list {
		notifications = append(notifications, s.CreateNotification(item))
	}

	if err := s.channel.WriteBatch(notifications); err != nil {
		return err
	}
	return nil
}

//CreateNotification creates notification from
func (s *Server) CreateNotification(event *model.Event) *queue.Notification {
	return &queue.Notification{
		Title:     event.Title,
		StartTime: event.StartTime,
		UserID:    event.UserID,
	}
}

//LogError writes error in log
func (s *Server) LogError(name string, err error) {
	s.logger.Error(
		"Scheduler server",
		zap.NamedError(name, err),
	)
}

package scheduler

import (
	"context"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"github.com/gzavodov/otus-go/calendar/app/queue"
	"go.uber.org/zap"
)

//NewServer creates new scheduler server
func NewServer(ctx context.Context, channel queue.NotificationChannel, repo repository.EventRepository, logger *zap.Logger) *Server {
	return &Server{ctx: ctx, channel: channel, repo: repo, logger: logger}
}

//Server event reminder service
type Server struct {
	ctx     context.Context
	channel queue.NotificationChannel
	repo    repository.EventRepository
	logger  *zap.Logger
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
		case <-time.After(1 * time.Minute):
			err = s.Check()
			if err != nil {
				s.LogError(ErrorCategoryRepository, err)
				return err
			}
		}
	}
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

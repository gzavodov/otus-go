package event

import (
	"context"
	"fmt"
	"time"

	"github.com/gzavodov/otus-go/calendar/factory/queuefactory"
	"github.com/gzavodov/otus-go/calendar/model"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
	"github.com/gzavodov/otus-go/calendar/pkg/queue"
	"github.com/gzavodov/otus-go/calendar/service/rpc"
	"github.com/gzavodov/otus-go/calendar/service/scheduler"
)

type Event struct {
	Repo                 *rpc.EventRepository
	Notifications        *NotificationReceiver
	NotificationChannel  chan *queue.Notification
	NotificationInterval int
	Client               *scheduler.Client
	InitialTime          time.Time
}

func NewEvent(ctx context.Context, conf *config.Configuration) (*Event, error) {
	now := time.Now().UTC()
	initialTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

	queueChannel, err := queuefactory.CreateQueueChannel(
		ctx,
		conf.AMPQTypeID,
		conf.AMPQName,
		conf.AMPQAddress,
	)

	if err != nil {
		return nil, err
	}

	notificationChannel := make(chan *queue.Notification)
	notifications := NewNotificationReceiver(notificationChannel)

	return &Event{
			Repo: rpc.NewEventRepository(ctx, conf.HTTPAddress),
			Client: scheduler.NewClient(
				ctx,
				queueChannel,
				notifications,
				nil,
			),
			Notifications:        notifications,
			NotificationChannel:  notificationChannel,
			NotificationInterval: conf.SchedulerCheckInterval,
			InitialTime:          initialTime,
		},
		nil
}

func (t *Event) Start() {
	go func(client *scheduler.Client) {
		if client != nil {
			client.Start()
		}
	}(t.Client)
}

func (t *Event) Stop() {
	go func(client *scheduler.Client) {
		if client != nil {
			client.Stop()
		}
	}(t.Client)
}

func (t *Event) CheckIfExist(userID int64, from time.Time, to time.Time) (bool, error) {
	events, err := t.Repo.ReadList(userID, from, to)
	if err != nil {
		return false, err
	}
	return len(events) > 0, nil
}

func (t *Event) EnsureExist(userID int, startAfter int, duration int) error {
	fromTime := t.InitialTime.Add(time.Duration(startAfter) * time.Minute)
	toTime := fromTime.Add(time.Duration(duration) * time.Minute)

	result, err := t.CheckIfExist(int64(userID), fromTime, toTime)
	if err != nil {
		return fmt.Errorf("function CheckIfExist failed with error: %w ", err)
	}

	if !result {
		return fmt.Errorf("calendar does not contain an event with start time after %d minutes", startAfter)
	}

	return nil
}

func (t *Event) WaitForNotification(title string) error {
	if t.Notifications.FindByTitle(title) != nil {
		return nil
	}

	timeout := t.NotificationInterval
	if timeout <= 0 {
		timeout = 60
	}
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		break
	case <-t.NotificationChannel:
		if t.Notifications.FindByTitle(title) != nil {
			return nil
		}
	}

	return fmt.Errorf("calendar does not contain an event '%s'", title)
}

func (t *Event) Create(m *model.Event) error {
	if err := t.Repo.Create(m); err != nil {
		return fmt.Errorf("failed to create event(%w)", err)
	}

	return nil
}

//GetDayList returns calendar events for specified date
func (t *Event) GetDayList(userID int64, date time.Time) ([]*model.Event, error) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)

	list, err := t.Repo.ReadList(userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve day list(%w)", err)
	}
	return list, nil
}

//GetWeekList returns calendar events for week that specified by date
func (t *Event) GetWeekList(userID int64, date time.Time) ([]*model.Event, error) {
	//week starts from monday
	dayIndex := int(date.Weekday())
	if dayIndex > 0 {
		dayIndex--
	} else {
		dayIndex = 6
	}
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1*dayIndex)
	to := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC).AddDate(0, 0, 6-dayIndex)

	list, err := t.Repo.ReadList(userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve week list(%w)", err)
	}
	return list, nil
}

//GetMonthList returns calendar events for month that specified by date
func (t *Event) GetMonthList(userID int64, date time.Time) ([]*model.Event, error) {
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 0, time.UTC)

	list, err := t.Repo.ReadList(userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve month list(%w)", err)
	}
	return list, nil
}

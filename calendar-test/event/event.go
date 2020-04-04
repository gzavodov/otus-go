package event

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gzavodov/otus-go/calendar/factory/queuefactory"
	"github.com/gzavodov/otus-go/calendar/model"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
	"github.com/gzavodov/otus-go/calendar/pkg/queue"
	"github.com/gzavodov/otus-go/calendar/service/rpc"
	"github.com/gzavodov/otus-go/calendar/service/scheduler"
)

//Event List Type
const (
	UnknownListType = 0
	DayListType     = 1
	WeekListType    = 2
	MonthListType   = 3
)

//NewEvent creates new event test according to configuration
func NewEvent(ctx context.Context, conf *config.Configuration) (*Event, error) {
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
			Repo: rpc.NewEventRepository(ctx, conf.GRPCAddress),
			Client: scheduler.NewClient(
				ctx,
				queueChannel,
				notifications,
				nil,
			),
			Notifications:        notifications,
			NotificationChannel:  notificationChannel,
			NotificationInterval: conf.SchedulerCheckInterval,
		},
		nil
}

//Event epresent facilities for testing calendar events
type Event struct {
	Repo                 *rpc.EventRepository
	Notifications        *NotificationReceiver
	NotificationChannel  chan *queue.Notification
	NotificationInterval int
	Client               *scheduler.Client
}

//Start run services and clients are requred for test process
func (t *Event) Start(outline interface{}) {
	go func(client *scheduler.Client) {
		if client != nil {
			client.Start()
		}
	}(t.Client)
}

//Stop halt services and clients are requred for test process
func (t *Event) Stop(outline interface{}, err error) {
	go func(client *scheduler.Client) {
		if client != nil {
			client.Stop()
		}
	}(t.Client)
}

//WaitForNotification waits for calendar event notification that specified by the title of event
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

//Create creates calendar event that specified by the model
func (t *Event) Create(m *model.Event) error {
	if err := t.Repo.Create(m); err != nil {
		return fmt.Errorf("failed to create event(%w)", err)
	}

	return nil
}

//GetDayList returns calendar events for specified the date
func (t *Event) GetDayList(userID int64, date time.Time) ([]*model.Event, error) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)

	list, err := t.Repo.ReadList(userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve day list(%w)", err)
	}
	return list, nil
}

//GetWeekList returns calendar events for week that specified by the date
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

//GetMonthList returns calendar events for month that specified by the date
func (t *Event) GetMonthList(userID int64, date time.Time) ([]*model.Event, error) {
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 0, time.UTC)

	list, err := t.Repo.ReadList(userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve month list(%w)", err)
	}
	return list, nil
}

//parseDateTime parses time from string by the RFC3339 format
//This function supports offsets from current time.
//Examples:
// "2020-03-02T12:00:00Z": time in the RFC3339 format
// "+30": current time + 30 minutes
// "-15": current time - 15 minutes
func (t *Event) parseDateTime(str string) (time.Time, error) {

	//NOW: current time
	if strings.EqualFold(str, "now") {
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC), nil
	}

	//[+-]/d+: offset from current time (negative or positive)
	rx, _ := regexp.Compile(`^([+-]?)(\d+)$`)
	matches := rx.FindAllStringSubmatch(str, -1)

	if len(matches) > 0 {
		match := matches[0]
		if len(match) > 2 {
			offset, err := strconv.Atoi(match[2])
			if err != nil {
				return time.Time{}, fmt.Errorf("could not parse time offset value (%w)", err)
			}

			now := time.Now().UTC()
			duration := time.Duration(offset) * time.Minute
			sign := match[1]
			if sign == "-" {
				duration *= -1
			}
			return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC).Add(duration), nil
		}
	}

	return time.Parse(time.RFC3339, str)
}

//parseModelFromTable parses event mmodel from data table
func (t *Event) parseModelFromTable(table *gherkin.DataTable, index int) (*model.Event, error) {
	dataRow := table.Rows[index]

	title := dataRow.Cells[0].Value

	startTime, err := t.parseDateTime(dataRow.Cells[1].Value)
	if err != nil {
		return nil, err
	}

	endTime, err := t.parseDateTime(dataRow.Cells[2].Value)
	if err != nil {
		return nil, err
	}

	notifyBefore, err := strconv.Atoi(dataRow.Cells[3].Value)
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(dataRow.Cells[4].Value)
	if err != nil {
		return nil, err
	}

	result := &model.Event{
		Title:        title,
		StartTime:    startTime,
		EndTime:      endTime,
		NotifyBefore: time.Duration(notifyBefore) * time.Minute,
		UserID:       int64(userID),
	}

	return result, nil
}

//findIndexByTitle searches event by the title in the list
func (t *Event) findIndexByTitle(list []*model.Event, title string) int {
	for i, eventModel := range list {
		if eventModel.Title == title {
			return i
		}
	}
	return -1
}

//CreateFirstFromTable creates event from the first data row of table
//The zero-index row treated as header. The data rows start from index 1
func (t *Event) CreateFirstFromTable(table *gherkin.DataTable) error {
	eventModel, err := t.parseModelFromTable(table, 1)
	if err != nil {
		return err
	}

	return t.Create(eventModel)
}

//CreateAllFromTable creates events from all the data rows of table
//The zero-index row treated as header. The data rows start from index 1
func (t *Event) CreateAllFromTable(table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		eventModel, err := t.parseModelFromTable(table, i)
		if err != nil {
			return err
		}

		err = t.Create(eventModel)
		if err != nil {
			return err
		}
	}
	return nil
}

//VerifyListByTable ensures that all the data rows from table present in schedule specified by type ID
func (t *Event) VerifyListByTable(table *gherkin.DataTable, listTypeID int) error {
	if listTypeID != DayListType && listTypeID != WeekListType && listTypeID != MonthListType {
		return fmt.Errorf("list type %d is not supported in current context", listTypeID)
	}

	var list []*model.Event
	var item *model.Event
	var err error

	for i := 1; i < len(table.Rows); i++ {
		item, err = t.parseModelFromTable(table, i)
		if err != nil {
			return err
		}

		switch listTypeID {
		case MonthListType:
			list, err = t.GetMonthList(item.UserID, item.StartTime)
		case WeekListType:
			list, err = t.GetWeekList(item.UserID, item.StartTime)
		case DayListType:
			list, err = t.GetDayList(item.UserID, item.StartTime)
		default:
			list = nil
			err = fmt.Errorf("list type %d is not supported in current context", listTypeID)
		}

		if err != nil {
			return err
		}

		if t.findIndexByTitle(list, item.Title) < 0 {
			return fmt.Errorf("could not find event with title '%s' in schedule", item.Title)
		}
	}
	return nil
}

//VerifyDayByTable ensures that all the data rows from table present in dayly schedule
func (t *Event) VerifyDayByTable(table *gherkin.DataTable) error {
	return t.VerifyListByTable(table, DayListType)
}

//VerifyWeekByTable ensures that all the data rows from table present in weekly schedule
func (t *Event) VerifyWeekByTable(table *gherkin.DataTable) error {
	return t.VerifyListByTable(table, WeekListType)
}

//VerifyMonthByTable ensures that all the data rows from table present in montly schedule
func (t *Event) VerifyMonthByTable(table *gherkin.DataTable) error {
	return t.VerifyListByTable(table, MonthListType)
}

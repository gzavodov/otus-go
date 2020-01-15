package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gzavodov/otus-go/calendar-test/event"
	"github.com/gzavodov/otus-go/calendar/model"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
)

type EventTest struct {
	Event *event.Event
}

func (t *EventTest) Start(interface{}) {
	t.Event.Start()
}

func (t *EventTest) Stop(interface{}, error) {
	t.Event.Stop()
}

func (t *EventTest) ParseDateTime(str string) (time.Time, error) {

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
func (t *EventTest) ParseModelFromTable(table *gherkin.DataTable, index int) (*model.Event, error) {
	dataRow := table.Rows[index]

	title := dataRow.Cells[0].Value

	startTime, err := t.ParseDateTime(dataRow.Cells[1].Value)
	if err != nil {
		return nil, err
	}

	endTime, err := t.ParseDateTime(dataRow.Cells[2].Value)
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

func (t *EventTest) CreateEvent(table *gherkin.DataTable) error {

	eventModel, err := t.ParseModelFromTable(table, 1)
	if err != nil {
		return err
	}

	return t.Event.Create(eventModel)
}

func (t *EventTest) CreateAllEvents(table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		eventModel, err := t.ParseModelFromTable(table, i)
		if err != nil {
			return err
		}

		err = t.Event.Create(eventModel)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *EventTest) FindEventIndexByTitle(list []*model.Event, title string) int {
	for i, eventModel := range list {
		if eventModel.Title == title {
			return i
		}
	}
	return -1
}

func (t *EventTest) EnsureEventNotification(title string) error {
	return t.Event.WaitForNotification(title)
}

func (t *EventTest) EnsureEventInDayList(table *gherkin.DataTable) error {
	eventModel, err := t.ParseModelFromTable(table, 1)
	if err != nil {
		return err
	}

	list, err := t.Event.GetDayList(eventModel.UserID, eventModel.StartTime)
	if err != nil {
		return err
	}

	if t.FindEventIndexByTitle(list, eventModel.Title) < 0 {
		return fmt.Errorf("could not find event with title '%s' in dayly schedule", eventModel.Title)
	}
	return nil
}

func (t *EventTest) EnsureEventsInWeekList(table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		eventModel, err := t.ParseModelFromTable(table, i)
		if err != nil {
			return err
		}

		list, err := t.Event.GetWeekList(eventModel.UserID, eventModel.StartTime)
		if err != nil {
			return err
		}

		if t.FindEventIndexByTitle(list, eventModel.Title) < 0 {
			return fmt.Errorf("could not find event with title '%s' in weekly schedule", eventModel.Title)
		}
	}

	return nil
}

func (t *EventTest) EnsureEventsInMonthList(table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		eventModel, err := t.ParseModelFromTable(table, i)
		if err != nil {
			return err
		}

		list, err := t.Event.GetMonthList(eventModel.UserID, eventModel.StartTime)
		if err != nil {
			return err
		}

		if t.FindEventIndexByTitle(list, eventModel.Title) < 0 {
			return fmt.Errorf("could not find event with title '%s' in weekly schedule", eventModel.Title)
		}
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	configFilePath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "./config/config.development.json"
	}

	configuration := &config.Configuration{}
	err := configuration.Load(
		*configFilePath,
		&config.Configuration{
			LogFilePath: "stderr",
			LogLevel:    "debug",
		},
	)
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	event, err := event.NewEvent(context.Background(), configuration)
	if err != nil {
		log.Fatalf("Could not create test: %v", err)
	}

	eventTest := EventTest{Event: event}

	s.BeforeScenario(eventTest.Start)

	s.Step(`^User creates today event:$`, eventTest.CreateEvent)
	s.Step(`^User receives notification with title "([^"]*)"$`, eventTest.EnsureEventNotification)

	s.Step(`^User creates day event:$`, eventTest.CreateEvent)
	s.Step(`^User\'s daily schedule contains an event that has been created:$`, eventTest.EnsureEventInDayList)

	s.Step(`^User creates events for week:$`, eventTest.CreateAllEvents)
	s.Step(`^User\'s weekly schedule contains all events that has been created:$`, eventTest.EnsureEventsInWeekList)

	s.Step(`^User creates events for month:$`, eventTest.CreateAllEvents)
	s.Step(`^User\'s monthly schedule contains all events that has been created:$`, eventTest.EnsureEventsInMonthList)

	s.AfterScenario(eventTest.Stop)
}

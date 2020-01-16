package main

import (
	"context"
	"flag"
	"log"

	"github.com/DATA-DOG/godog"
	"github.com/gzavodov/otus-go/calendar-test/event"
	"github.com/gzavodov/otus-go/calendar/pkg/config"
)

//FeatureContext implements godog library entry point
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
		log.Fatalf("failed to load configuration: %v", err)
	}

	eventTest, err := event.NewEvent(context.Background(), configuration)
	if err != nil {
		log.Fatalf("faile to create test: %v", err)
	}

	s.BeforeScenario(eventTest.Start)

	s.Step(`^User creates today event:$`, eventTest.CreateFirstFromTable)
	s.Step(`^User receives notification with title "([^"]*)"$`, eventTest.WaitForNotification)

	s.Step(`^User creates day event:$`, eventTest.CreateAllFromTable)
	s.Step(`^User\'s daily schedule contains an event that has been created:$`, eventTest.VerifyDayByTable)

	s.Step(`^User creates events for week:$`, eventTest.CreateAllFromTable)
	s.Step(`^User\'s weekly schedule contains all events that has been created:$`, eventTest.VerifyWeekByTable)

	s.Step(`^User creates events for month:$`, eventTest.CreateAllFromTable)
	s.Step(`^User\'s monthly schedule contains all events that has been created:$`, eventTest.VerifyMonthByTable)

	s.AfterScenario(eventTest.Stop)
}

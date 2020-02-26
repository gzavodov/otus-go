package main

import (
	"context"
	"log"

	"github.com/cucumber/godog"
	"github.com/gzavodov/otus-go/banner-rotation/config"
)

//FeatureContext implements godog library entry point
func FeatureContext(s *godog.Suite) {
	configuration := &config.Configuration{}
	err := configuration.Load(
		"",
		&config.Configuration{
			LogFilePath: "stderr",
			LogLevel:    "debug",
		},
	)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	t, err := NewFeatureTest(context.Background(), configuration)
	if err != nil {
		log.Fatalf("failed to create test: %v", err)
	}

	s.BeforeFeature(t.Start)

	s.Step(`^Client creates following banner social groups:$`, t.CreateSocialGroupsFromTable)
	s.Step(`^Recently created banner social groups are available for using:$`, t.VerifySocialGroupsFromTable)

	s.Step(`^Client creates following banner slots:$`, t.CreateSlotsFromTable)
	s.Step(`^Recently created banner slots are available for using:$`, t.VerifySlotsFromTable)

	s.Step(`^Client creates following banners and bind to specified slots:$`, t.CreateBannersFromTable)
	s.Step(`^Recently created banner slots are available for using and bound to appropriate slots:$`, t.VerifyBannersFromTable)

	s.Step(`^Client makes query about banner show for following slots and social groups:$`, t.ChooseBanner)
	s.Step(`^Client receives notification about banner show$`, t.WaitForBannerChoiseNotification)

	s.Step(`^Client registers banner click event for banner selected on previous step:$`, t.RegisterBannerClick)
	s.Step(`^Client receives notification about banner click$`, t.WaitForBannerClickNotification)

	s.AfterFeature(t.Stop)
}

package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	exitCode := godog.RunWithOptions(
		"Banner Rotation Integration Test",
		func(s *godog.Suite) { FeatureContext(s) },
		godog.Options{
			Format:      "progress",
			Paths:       []string{"features"},
			Randomize:   0,
			Concurrency: 2,
		},
	)

	if code := m.Run(); code > exitCode {
		exitCode = code
	}

	os.Exit(exitCode)
}

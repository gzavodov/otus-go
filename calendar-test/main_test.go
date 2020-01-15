package main

import (
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
)

func TestMain(m *testing.M) {
	exitCode := godog.RunWithOptions(
		"Calendar Integration Test",
		func(s *godog.Suite) { FeatureContext(s) },
		godog.Options{
			Format:    "progress",
			Paths:     []string{"features"},
			Randomize: 0,
		},
	)

	if code := m.Run(); code > exitCode {
		exitCode = code
	}

	os.Exit(exitCode)
}

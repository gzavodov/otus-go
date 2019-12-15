package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"go.uber.org/zap/zapcore"
)

//Configuration Web server configuration
type Configuration struct {
	HTTPAddress string `json:"http_listen"`

	LogFilePath string `json:"log_file"`
	LogLevel    string `json:"log_level"`

	EndpointServiceTypeID int `json:"endpoint_service_type_id"`

	EventRepositoryTypeID int    `json:"event_repository_type_id"`
	EventRepositoryDSN    string `json:"event_repository_dsn"`

	AMPQTypeID  int    `json:"ampq_type_id"`
	AMPQName    string `json:"ampq_name"`
	AMPQAddress string `json:"ampq_address"`
}

//Load read configuration from file anf from os environment variables
func (c *Configuration) Load(filePath string, defaultVal *Configuration) error {
	if len(filePath) > 0 {
		err := c.LoadFromFile(filePath)
		if err != nil {
			return err
		}
	}

	c.LoadFromEvironment()

	if defaultVal != nil {
		if c.HTTPAddress == "" {
			c.HTTPAddress = defaultVal.HTTPAddress
		}

		if c.LogFilePath == "" {
			c.LogFilePath = defaultVal.LogFilePath
		}

		if c.LogLevel == "" {
			c.LogLevel = defaultVal.LogLevel
		}

		if c.EndpointServiceTypeID <= 0 {
			c.EndpointServiceTypeID = defaultVal.EndpointServiceTypeID
		}

		if c.EventRepositoryTypeID <= 0 {
			c.EventRepositoryTypeID = defaultVal.EventRepositoryTypeID
		}

		if c.EventRepositoryDSN == "" {
			c.EventRepositoryDSN = defaultVal.EventRepositoryDSN
		}

		if c.AMPQTypeID <= 0 {
			c.AMPQTypeID = defaultVal.AMPQTypeID
		}

		if c.AMPQName == "" {
			c.AMPQName = defaultVal.AMPQName
		}

		if c.AMPQAddress == "" {
			c.AMPQAddress = defaultVal.AMPQAddress
		}
	}

	return nil
}

//LoadFromFile read configuration from file
func (c *Configuration) LoadFromFile(filePath string) error {
	configFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Could not read configuration file: %w", err)
	}

	if json.Unmarshal(configFile, c) != nil {
		return fmt.Errorf("Could not internalize configuration file data: %w", err)
	}

	return nil
}

//LoadFromEvironment read configuration from environment variables
func (c *Configuration) LoadFromEvironment() error {
	if s, ok := os.LookupEnv("CALENDAR_REPOSITORY_TYPE"); ok {
		result, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("Could not parse CALENDAR_REPOSITORY_TYPE variable: %w", err)
		}
		c.EndpointServiceTypeID = int(result)
	}

	if s, ok := os.LookupEnv("CALENDAR_REPOSITORY_DSN"); ok {
		c.EventRepositoryDSN = s
	}

	return nil
}

//ParseZapLogLevel Parse Zap LogLevel from strig
func (c *Configuration) ParseZapLogLevel() (zapcore.Level, error) {
	var l zapcore.Level
	err := l.UnmarshalText([]byte(c.LogLevel))
	return l, err
}

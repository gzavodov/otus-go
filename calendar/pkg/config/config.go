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
	HealthcheckHTTPAddress string `json:"healthcheck_http_listen"`
	HTTPAddress            string `json:"http_listen"`
	GRPCAddress            string `json:"grpc_listen"`

	LogFilePath string `json:"log_file"`
	LogLevel    string `json:"log_level"`

	ServiceTypeID int `json:"service_type_id"`

	EventRepositoryTypeID int    `json:"event_repository_type_id"`
	EventRepositoryDSN    string `json:"event_repository_dsn"`

	AMPQTypeID  int    `json:"ampq_type_id"`
	AMPQName    string `json:"ampq_name"`
	AMPQAddress string `json:"ampq_address"`

	SchedulerCheckInterval int `json:"scheduler_check_interval"`
}

//Load read configuration from file anf from os environment variables
func (c *Configuration) Load(filePath string, defaultVal *Configuration) error {
	if len(filePath) > 0 {
		if err := c.LoadFromFile(filePath); err != nil {
			return err
		}
	}

	if err := c.LoadFromEvironment(); err != nil {
		return err
	}

	if defaultVal != nil {
		if c.HTTPAddress == "" {
			c.HTTPAddress = defaultVal.HTTPAddress
		}

		if c.GRPCAddress == "" {
			c.GRPCAddress = defaultVal.GRPCAddress
		}

		if c.HealthcheckHTTPAddress == "" {
			c.HealthcheckHTTPAddress = defaultVal.HealthcheckHTTPAddress
		}

		if c.LogFilePath == "" {
			c.LogFilePath = defaultVal.LogFilePath
		}

		if c.LogLevel == "" {
			c.LogLevel = defaultVal.LogLevel
		}

		if c.ServiceTypeID <= 0 {
			c.ServiceTypeID = defaultVal.ServiceTypeID
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

		if c.SchedulerCheckInterval <= 0 {
			c.SchedulerCheckInterval = defaultVal.SchedulerCheckInterval
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
	//HTTP End Point Address
	if s, ok := os.LookupEnv("CALENDAR_HTTP_ADDRESS"); ok {
		c.HTTPAddress = s
	}

	//GRPC End Point Address
	if s, ok := os.LookupEnv("CALENDAR_GRPC_ADDRESS"); ok {
		c.GRPCAddress = s
	}

	//Monitoring End Point Address
	if s, ok := os.LookupEnv("CALENDAR_HEALTH_CHECK_HTTP_ADDRESS"); ok {
		c.HealthcheckHTTPAddress = s
	}

	//Repository
	if s, ok := os.LookupEnv("CALENDAR_REPOSITORY_TYPE"); ok {
		result, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("Could not parse CALENDAR_REPOSITORY_TYPE variable: %w", err)
		}
		c.EventRepositoryTypeID = int(result)
	}

	if s, ok := os.LookupEnv("CALENDAR_REPOSITORY_DSN"); ok {
		c.EventRepositoryDSN = s
	}

	//Service
	if s, ok := os.LookupEnv("CALENDAR_SERVICE_TYPE"); ok {
		result, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("Could not parse CALENDAR_SERVICE_TYPE variable: %w", err)
		}
		c.ServiceTypeID = int(result)
	}

	//AMPQ
	if s, ok := os.LookupEnv("AMPQ_TYPE_ID"); ok {
		result, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("Could not parse AMPQ_TYPE_ID variable: %w", err)
		}
		c.AMPQTypeID = int(result)
	}

	if s, ok := os.LookupEnv("AMPQ_NAME"); ok {
		c.AMPQName = s
	}

	if s, ok := os.LookupEnv("AMPQ_ADDRESS"); ok {
		c.AMPQAddress = s
	}

	return nil
}

//ParseZapLogLevel Parse Zap LogLevel from strig
func (c *Configuration) ParseZapLogLevel() (zapcore.Level, error) {
	var l zapcore.Level
	err := l.UnmarshalText([]byte(c.LogLevel))
	return l, err
}

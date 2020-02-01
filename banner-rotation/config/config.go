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

	RepositoryDSN string `json:"repository_dsn"`

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

		if c.RepositoryDSN == "" {
			c.RepositoryDSN = defaultVal.RepositoryDSN
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
	//End Point IP-Address
	if s, ok := os.LookupEnv("BANNER_ROTATION_HTTP_ADDRESS"); ok {
		c.HTTPAddress = s
	}

	//Repository
	if s, ok := os.LookupEnv("BANNER_ROTATION_REPOSITORY_DSN"); ok {
		c.RepositoryDSN = s
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

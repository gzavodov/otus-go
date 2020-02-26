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

	AlgorithmTypeID int `json:"algorithm_type_id"`

	RepositoryDSN string `json:"repository_dsn"`

	AMPQName    string `json:"ampq_name"`
	AMPQAddress string `json:"ampq_address"`
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

	c.setIfEmpty(defaultVal)

	return nil
}

func (c *Configuration) setIfEmpty(source *Configuration) {
	if source == nil {
		return
	}

	if c.HTTPAddress == "" {
		c.HTTPAddress = source.HTTPAddress
	}

	if c.LogFilePath == "" {
		c.LogFilePath = source.LogFilePath
	}

	if c.LogLevel == "" {
		c.LogLevel = source.LogLevel
	}

	if c.RepositoryDSN == "" {
		c.RepositoryDSN = source.RepositoryDSN
	}

	if c.AMPQName == "" {
		c.AMPQName = source.AMPQName
	}

	if c.AMPQAddress == "" {
		c.AMPQAddress = source.AMPQAddress
	}
}

//LoadFromFile read configuration from file
func (c *Configuration) LoadFromFile(filePath string) error {
	configFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read configuration file: %w", err)
	}

	if json.Unmarshal(configFile, c) != nil {
		return fmt.Errorf("could not internalize configuration file data: %w", err)
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

	//Algorithm
	if s, ok := os.LookupEnv("ALGORITHM_TYPE_ID"); ok {
		result, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("could not parse ALGORITHM_TYPE_ID variable: %w", err)
		}
		c.AlgorithmTypeID = int(result)
	}

	//AMPQ
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

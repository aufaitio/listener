package app

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/viper"
)

// Config stores the application-wide configurations
var Config AppConfig

// AppConfig configuration necessary for the listener API
type AppConfig struct {
	DB        dbConfig
	ErrorFile string
	Port      int32
}

// DBConfig Config representing database info.
type dbConfig struct {
	Host     string
	Name     string
	Password string
	Port     int32
	Username string
}

// Validate validates AppConfig, currently unused but keeping around in case it is needed
func (config AppConfig) Validate() error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.DB, validation.Required),
	)
}

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("app")
	v.AutomaticEnv()
	v.SetDefault("ErrorFile", "config/errors.yaml")
	v.SetDefault("Port", 8080)
	v.SetDefault("DB", dbConfig{Host: "localhost", Port: 27017, Name: "aufait"})

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}

	if err := v.Unmarshal(&Config); err != nil {
		return err
	}

	return nil
}

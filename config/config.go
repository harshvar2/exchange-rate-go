package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type AppConfig struct {
	ExchangeAPIKey string
}

var Config AppConfig

func InitConfig() {
	// Get the absolute path to config.yml
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get caller information")
	}
	configFilePath := filepath.Join(filepath.Dir(filename), "../config.yml")

	// Set config file and allow env override
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config: %v", err)
		}
	}

	Config.ExchangeAPIKey = viper.GetString("EXCHANGE_API_KEY")
	if Config.ExchangeAPIKey == "" {
		Config.ExchangeAPIKey = os.Getenv("EXCHANGE_API_KEY")
	}
	if Config.ExchangeAPIKey == "" {
		log.Fatal("EXCHANGE_API_KEY must be set in config.yml or environment")
	}
}

package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App   AppConfig
	Kafka KafkaConfig
}

type AppConfig struct {
	Environment    string
	Id             string
	PathToWatch    string
	ErrorParsePath string
	Workers        int
}

type KafkaConfig struct {
	Brokers        []string
	TopicsConsumer []string
	TopicsProducer []string
	TopicsError    []string
	GroupID        string
	MaxRetry       int
}

type configValidator func(*Config) error

var validators = []configValidator{
	func(cfg *Config) error {
		return validateEnvironment(cfg.App.Environment)
	},
}

var internalConfig *Config

func InitConfig(configPath string) {
	viper.SetConfigName("config")
	if configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath(".")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("app.environment", "prod")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed: %v", err)
	}

	var configMap map[string]interface{}
	if err := viper.Unmarshal(&configMap); err != nil {
		log.Fatalf("unmarshal to map failed: %v", err)
	}

	if err := mapstructure.Decode(configMap, &internalConfig); err != nil {
		log.Fatalf("decode map to struct failed: %v", err)
	}

	for _, validator := range validators {
		if err := validator(internalConfig); err != nil {
			log.Fatalf("configuration error: %v", err)
		}
	}
}

func Kafka() *KafkaConfig {
	return &internalConfig.Kafka
}

func App() *AppConfig {
	return &internalConfig.App
}

// Validators
func validateEnvironment(env string) error {
	validEnvironments := map[string]bool{"dev": true, "prod": true}
	if _, isValid := validEnvironments[env]; !isValid {
		return fmt.Errorf("invalid environment '%s': the environment must be either 'dev' or 'prod'", env)
	}
	return nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port number %d: port must be between 1 and 65535", port)
	}
	return nil
}

package config

import (
	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

// Config sets the configuration for the gRPC service
type Config struct {
	Env             string   `env:"ENV" envDefault:"dev"`
	ConsumerGroupID string   `env:"CONSUMER_GROUP_ID"`
	Port            uint16   `env:"PORT" envDefault:"9090"`
	Seeds           []string `env:"BROKERS" envSeparator:","`
	Topics          []string `env:"TOPICS" envSeparator:","`
}

var log *zap.SugaredLogger

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log = logger.Sugar()
}

func New() *Config {
	var config = new(Config)

	if err := env.Parse(config); err != nil {
		log.Fatalf("error parsing config, %v", err)
	}

	return config
}

func (c *Config) DefaultProducerTopic() string {
	return c.Topics[0]
}

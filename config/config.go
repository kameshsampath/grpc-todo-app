package config

import (
	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

// Config sets the configuration for the gRPC server
type Config struct {
	Env             string   `env:"ENV" envDefault:"dev"`
	ConsumerGroupID string   `env:"CONSUMER_GROUP_ID"`
	Port            uint16   `env:"PORT" envDefault:"9090"`
	Seeds           []string `env:"BROKERS" envSeparator:"," envDefault:"localhost:19092"`
	Topics          []string `env:"TOPICS" envSeparator:","`
	SchemaRegistry  string   `env:"SCHEMA_REGISTRY" envDefault:"localhost:18081"`
	SchemaURL       string   `env:"SCHEMA_URL" envDefault:"https://raw.githubusercontent.com/kameshsampath/demo-protos/main/todo/todo.proto"`
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

// DefaultProducerTopic gets the default topic that will be used as the producer topic
func (c *Config) DefaultProducerTopic() string {
	return c.Topics[0]
}

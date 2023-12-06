package main

import (
	"github.com/kameshsampath/todo-app/config"
	"github.com/kameshsampath/todo-app/internal/impl"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func main() {
	//Jai Guru
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log := logger.Sugar()

	config := config.New()

	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.Seeds...),
		kgo.ConsumeTopics(config.Topics...),
		kgo.DefaultProduceTopic(config.DefaultProducerTopic()),
		kgo.ConsumerGroup(config.ConsumerGroupID),
		kgo.AllowAutoTopicCreation(),
	)

	if err != nil {
		log.Fatal(err)
	}

	server := impl.New(client, config)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

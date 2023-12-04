package main

import (
	"github.com/kameshsampath/todo-app/config"
	"github.com/kameshsampath/todo-app/internal/adapters/grpc"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func main() {
	//Jai Guru
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log := logger.Sugar()

	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.GetBrokers()...),
		kgo.ConsumeTopics(config.GetTopics()...),
		kgo.DefaultProduceTopic(config.GetDefaultProducerTopic()),
		kgo.ConsumerGroup(config.GetConsumerGroup()),
	)

	if err != nil {
		log.Fatal(err)
	}

	adapter := grpc.NewAdapter(client, config.GetPort())
	if err := adapter.Run(); err != nil {
		log.Fatal(err)
	}
}

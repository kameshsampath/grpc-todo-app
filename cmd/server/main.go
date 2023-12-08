package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kameshsampath/demo-protos/golang/todo"
	"github.com/kameshsampath/todo-app/config"
	"github.com/kameshsampath/todo-app/impl"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type serviceSchemas struct {
	id        int
	typeValue any
	index     int
}

var (
	protoMarshallFn = func(a any) ([]byte, error) {
		return proto.Marshal(a.(proto.Message))
	}
	protoUnMarshallFn = func(b []byte, a any) error {
		return proto.Unmarshal(b, a.(proto.Message))
	}
	log *zap.SugaredLogger
)

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log = logger.Sugar()
}

func main() {

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

	ss, err := createSchema(*config)
	if err != nil {
		log.Fatal(err)
	}
	// Add all schema type that needs to be registered with srede
	// and used when encoding and decoding proto messages
	serde := registerSchemas(
		[]serviceSchemas{
			{
				id:        ss.ID,
				index:     0,
				typeValue: &todo.Task{},
			},
		},
	)

	server := impl.New(
		impl.WithClient(client),
		impl.WithConfig(config),
		impl.WithSerde(serde),
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

// createSchema creates(registers) the schemas with the SchemaRegistry
func createSchema(config config.Config) (*sr.SubjectSchema, error) {
	rcl, err := sr.NewClient(sr.URLs(config.SchemaRegistry))
	if err != nil {
		return nil, err
	}

	res, err := http.Get(config.SchemaURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		sb, err := io.ReadAll(res.Body)
		log.Debugf("%s", string(sb))
		if err != nil {
			return nil, err
		}

		// if err := rcl.DeleteSchema(context.Background(), fmt.Sprintf("%s-value", config.DefaultProducerTopic()), -1, sr.HardDelete); err != nil {
		// 	log.Fatal(err)
		// }
		// os.Exit(1)

		ss, err := rcl.CreateSchema(context.Background(), fmt.Sprintf("%s-value", config.DefaultProducerTopic()), sr.Schema{
			Type:   sr.TypeProtobuf,
			Schema: string(sb),
		})
		if err != nil {
			return nil, err
		}

		return &ss, nil
	}
	return nil, fmt.Errorf("unable to read schema from url, '%s'", config.SchemaURL)
}

// registerSchemas registers all the schema and its corresponding encoding/decoding function
func registerSchemas(serviceSchemas []serviceSchemas) *sr.Serde {
	serde := new(sr.Serde)
	for _, schema := range serviceSchemas {
		serde.Register(
			schema.id,
			schema.typeValue,
			sr.DecodeFn(protoUnMarshallFn),
			sr.EncodeFn(protoMarshallFn),
			sr.Index(schema.index),
		)
	}
	return serde
}

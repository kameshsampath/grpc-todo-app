package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/kameshsampath/demo-protos/golang/todo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	//Jai Guru

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log := logger.Sugar()

	conn, err := grpc.Dial(os.Getenv("SERVICE_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := todo.NewTodoClient(conn)
	stream, err := client.TodoList(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		tr, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.TodoList failed: %v", err)
		}
		log.Infow("Task",
			"Title", tr.Task.Title,
			"Description", tr.Task.Description,
			"Completed", tr.Task.Completed,
			"Last Updated", tr.Task.LastUpdated.AsTime().Format(time.RFC850),
			"Partition", tr.Partition,
			"Offset", tr.Offset,
		)
	}
}

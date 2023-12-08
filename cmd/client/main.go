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

		if errs := tr.GetErrors(); errs != nil {
			log.Errorln("Errors fetching:")
			for _, e := range errs.Error {
				log.Errorw("Error Details",
					"Topic", e.Topic,
					"Partition", e.Partition,
					"Error", e.Message,
				)
			}
		} else {
			todo := tr.GetTodo()
			log.Infow("Task",
				"Title", todo.Task.Title,
				"Description", todo.Task.Description,
				"Completed", todo.Task.Completed,
				"Last Updated", todo.Task.LastUpdated.AsTime().Format(time.RFC850),
				"Partition", todo.Partition,
				"Offset", todo.Offset,
			)
		}
	}
}

package grpc

import (
	"context"
	"time"

	"github.com/kameshsampath/demo-protos/golang/todo"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type result struct {
	record *kgo.Record
	errors []kgo.FetchError
}

var log *zap.SugaredLogger

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log = logger.Sugar()
}

// AddTodo implements todo.TodoServer.
func (a *Adapter) AddTodo(ctx context.Context, req *todo.TodoAddRequest) (*todo.TodoResponse, error) {
	task := &todo.Task{
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Completed:   req.Task.Completed,
	}
	return marshallAndSend(ctx, a.client, task)
}

// UpdateTodo implements todo.TodoServer.
func (a *Adapter) UpdateTodo(ctx context.Context, req *todo.UpdateTodoStatusRequest) (*todo.TodoResponse, error) {
	task := &todo.Task{
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Completed:   req.Task.Completed,
	}
	return marshallAndSend(ctx, a.client, task)
}

// TodoList implements todo.TodoServer.
func (a *Adapter) TodoList(empty *emptypb.Empty, stream todo.Todo_TodoListServer) error {
	ch := make(chan result)
	go func() {
		poll(a.client, ch)
	}()

	for {
		select {
		case r := <-ch:
			{
				if errs := r.errors; len(errs) > 0 {
					log.Errorln("Errors")
					//TODO send back errors
					for _, err := range errs {
						log.Errorw("Error Details",
							"Topic", err.Topic,
							"Partition", err.Partition,
							"Error", err.Err,
						)
					}
				}
				b := r.record.Value
				task := new(todo.Task)
				if err := proto.Unmarshal(b, task); err != nil {
					//TODO send back errors
					log.Error(err)
				} else {
					stream.Send(&todo.TodoResponse{
						Task:      task,
						Partition: r.record.Partition,
						Offset:    r.record.Offset,
					})
				}
			}
		}
	}
}

// marshallAndSend sends the task record to backend
func marshallAndSend(ctx context.Context, client *kgo.Client, task *todo.Task) (*todo.TodoResponse, error) {
	tctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	out, err := proto.Marshal(task)
	if err != nil {
		return nil, err
	}
	r := &kgo.Record{
		Key:   []byte(task.Title),
		Value: out,
	}
	if err := client.ProduceSync(tctx, r).FirstErr(); err != nil {
		return nil, err
	}
	return &todo.TodoResponse{
		Task:      task,
		Partition: r.Partition,
		Offset:    r.Offset,
	}, nil
}

// poll fetches the record from the backend and adds that the channel
func poll(client *kgo.Client, ch chan result) {
	log.Debugln("Started to poll topic")
	//Consumer
	for {
		fetches := client.PollFetches(context.Background())
		if errs := fetches.Errors(); len(errs) > 0 {
			ch <- result{
				errors: errs,
			}
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, r := range p.Records {
				ch <- result{
					record: r,
				}
			}
		})
	}
}

var _ todo.TodoServer = (*Adapter)(nil)

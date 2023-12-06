package impl

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
func (s *Server) AddTodo(ctx context.Context, req *todo.TodoAddRequest) (*todo.TodoResponse, error) {
	task := &todo.Task{
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Completed:   req.Task.Completed,
	}
	return marshallAndSend(ctx, s.client, task)
}

// TodoList implements todo.TodoServer.
func (s *Server) TodoList(empty *emptypb.Empty, stream todo.Todo_TodoListServer) error {
	ch := make(chan result)
	go func() {
		poll(s.client, ch)
	}()

	for {
		select {
		case r := <-ch:
			{
				if errs := r.errors; len(errs) > 0 {
					var errors = make([]*todo.Error, len(errs))
					for _, err := range errs {
						log.Debugf("Error Details",
							"Topic", err.Topic,
							"Partition", err.Partition,
							"Error", err.Err,
						)
						errors = append(errors, &todo.Error{
							Topic:     err.Topic,
							Partition: err.Partition,
							Message:   err.Err.Error(),
						})
					}
					stream.Send(&todo.TodoListResponse{
						Response: &todo.TodoListResponse_Errors{Errors: &todo.Errors{
							Error: errors,
						}},
					})
				}
				b := r.record.Value
				task := new(todo.Task)
				if err := proto.Unmarshal(b, task); err != nil {
					//Skip Sending invalid data, just log the error
					log.Errorw("Error marshalling task",
						"Data", string(b),
						"Error", err.Error())
				} else {
					stream.Send(&todo.TodoListResponse{
						Response: &todo.TodoListResponse_Todo{
							Todo: &todo.TodoResponse{
								Task:      task,
								Partition: r.record.Partition,
								Offset:    r.record.Offset,
							},
						},
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

var _ todo.TodoServer = (*Server)(nil)

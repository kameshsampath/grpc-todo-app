package grpc

import (
	"github.com/kameshsampath/demo-protos/golang/todo"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Adapter struct {
	client *kgo.Client
	port   int
	todo.UnimplementedTodoServer
}

// NewAdapter returns the adapter returns the Adapter to call gRPC services
func NewAdapter(client *kgo.Client, port int) *Adapter {
	return &Adapter{
		client: client,
		port:   port,
	}
}

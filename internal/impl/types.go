package impl

import (
	"github.com/kameshsampath/demo-protos/golang/todo"
	"github.com/kameshsampath/todo-app/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Server struct {
	client *kgo.Client
	config *config.Config
	todo.UnimplementedTodoServer
}

// Creates a new instance of the server with given configuration
func New(client *kgo.Client, config *config.Config) *Server {
	return &Server{
		client: client,
		config: config,
	}
}

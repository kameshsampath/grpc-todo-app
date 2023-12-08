package impl

import (
	"github.com/kameshsampath/demo-protos/golang/todo"
	"github.com/kameshsampath/todo-app/config"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

type Server struct {
	client *kgo.Client
	config *config.Config
	serde  *sr.Serde
	todo.UnimplementedTodoServer
}

type ServerOpt func(*Server)

// Creates a new instance of the server with given configuration
func New(options ...ServerOpt) *Server {
	s := &Server{}
	for _, opt := range options {
		opt(s)
	}

	return s
}

func WithClient(client *kgo.Client) ServerOpt {
	return func(s *Server) {
		s.client = client
	}
}

func WithSerde(serde *sr.Serde) ServerOpt {
	return func(s *Server) {
		s.serde = serde
	}
}

func WithConfig(config *config.Config) ServerOpt {
	return func(s *Server) {
		s.config = config
	}
}

package impl

import (
	"fmt"
	"net"

	"github.com/kameshsampath/demo-protos/golang/todo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run runs the gRPC server
func (s *Server) Run() error {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log := logger.Sugar()

	config := s.config
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return fmt.Errorf("error starting server with port %d,%v", config.Port, err)
	}
	server := grpc.NewServer()
	todo.RegisterTodoServer(server, s)
	log.Infof("Server started on port %d", config.Port)
	// required for grpcurl
	if config.Env == "dev" {
		reflection.Register(server)
	}

	if err := server.Serve(listen); err != nil {
		return fmt.Errorf("error starting server,%v", err)
	}
	return nil
}

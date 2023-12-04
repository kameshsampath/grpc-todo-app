package grpc

import (
	"fmt"
	"net"

	"github.com/kameshsampath/demo-protos/golang/todo"
	"github.com/kameshsampath/todo-app/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (a *Adapter) Run() error {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log := logger.Sugar()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("error starting server with port %d,%v", a.port, err)
	}
	server := grpc.NewServer()
	todo.RegisterTodoServer(server, a)
	// required for grpcurl
	if config.GetEnv() == "dev" {
		reflection.Register(server)
	}
	if err := server.Serve(listen); err != nil {
		return fmt.Errorf("error starting server,%v", err)
	}
	log.Infof("Server started on port %d", a.port)

	return nil
}

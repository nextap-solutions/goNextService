package components

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type GrpcComponent struct {
	grocServer *grpc.Server

	lis net.Listener

	exitChan chan (bool)
}

func NewGrpcComponent(server *grpc.Server, lis net.Listener) *GrpcComponent {
	component := GrpcComponent{
		grocServer: server,
		lis:        lis,
		exitChan:   make(chan bool),
	}
	return &component
}

func (gc *GrpcComponent) Startup() error {
	return nil
}

func (gc *GrpcComponent) Run() error {
	if gc.grocServer != nil {
		return gc.grocServer.Serve(gc.lis)
	}
	return nil
}

func (gc *GrpcComponent) Close(ctx context.Context) error {
	if gc.grocServer != nil {
		gc.grocServer.GracefulStop()
	}
	return nil
}

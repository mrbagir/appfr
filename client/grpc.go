package client

import (
	"github.com/mrbagir/appfr/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type app interface {
	Logger() logging.Logger
	RegisterClientConn(conn *grpc.ClientConn)
}

func NewGRPCClient[T any](app app, target string, newClient func(cc grpc.ClientConnInterface) T) T {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(target, options...)
	if err != nil {
		app.Logger().Fatalf("failed to create gRPC client connection to %s: %v", target, err)
	}

	app.Logger().Infof("gRPC client connected to %s", target)
	app.RegisterClientConn(conn)
	return newClient(conn)
}

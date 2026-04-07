package app

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient[T any](app *App, target string, newClient func(cc grpc.ClientConnInterface) T) T {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.logger.Errorf("failed to create gRPC client connection: %v", err)
		return newClient(nil)
	}
	app.logger.Infof("gRPC client connected to %s", target)
	return newClient(conn)
}

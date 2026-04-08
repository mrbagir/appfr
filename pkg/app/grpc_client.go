package app

import (
	"github.com/mrbagir/qcash-appcore/pkg/logging"
	"google.golang.org/grpc"
)

type grpcClient struct {
	address string
	conn    *grpc.ClientConn
}

type grpcClients struct {
	clients []grpcClient
	logger  logging.Logger
}

func (a *App) RegisterClientConn(conn *grpc.ClientConn) {
	if a.grpcClients == nil {
		a.grpcClients = &grpcClients{
			logger: a.logger,
		}
	}

	a.grpcClients.clients = append(a.grpcClients.clients, grpcClient{
		address: conn.Target(),
		conn:    conn,
	})
}

func (a *grpcClients) shutdown() {
	for _, c := range a.clients {
		if c.conn != nil {
			if err := c.conn.Close(); err != nil {
				a.logger.Errorf("failed to close gRPC client connection: %v", err)
			}
		}
	}
}

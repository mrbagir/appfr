package app

import (
	"github.com/mrbagir/qcash-appcore/pkg/config"
	"google.golang.org/grpc"
)

type grpcServer struct {
	server             *grpc.Server
	interceptors       []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	options            []grpc.ServerOption
	port               int
	config             config.Config
}

func (a *App) RegisterService(desc *grpc.ServiceDesc, impl any) {
	a.grpcServer.RegisterService(desc, impl)
}

func newGRPCServer(cfg config.Config) *grpcServer {

	return &grpcServer{
		server: grpc.NewServer(),
	}
}

func (g *grpcServer) Run() {
	if g.server == nil {
		g.server = grpc.NewServer(g.options...)
	}
}

func (g *grpcServer) RegisterService(desc *grpc.ServiceDesc, impl any) {
	g.server.RegisterService(desc, impl)
}

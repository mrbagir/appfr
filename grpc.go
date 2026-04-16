package appfr

import (
	"context"
	"net"
	"strconv"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"

	grpc_fr "github.com/mrbagir/appfr/grpc"
	"github.com/mrbagir/appfr/logging"
)

type grpcConfig struct {
	Port int `env:"GRPC_PORT" envDefault:"9000"`
}

type grpcServer struct {
	server             *grpc.Server
	interceptors       []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	options            []grpc.ServerOption
	config             grpcConfig
	logger             logging.Logger
}

func (a *App) RegisterService(desc *grpc.ServiceDesc, impl any) {
	if !a.grpcRegistered {
		a.grpcServer.createServer()
	}

	a.grpcServer.server.RegisterService(desc, impl)
	a.grpcRegistered = true
}

func newGRPCServer(logger logging.Logger, cfg config) *grpcServer {
	middleware := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
		grpc_fr.ObservabilityInterceptor(logger),
	}
	streamMiddleware := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(),
	}

	return &grpcServer{
		server:             grpc.NewServer(),
		interceptors:       middleware,
		streamInterceptors: streamMiddleware,
		config:             cfg.GrpcConfig,
		logger:             logger,
	}
}

func (g *grpcServer) createServer() {
	interceptorOption := grpc.ChainUnaryInterceptor(g.interceptors...)
	streamOpt := grpc.ChainStreamInterceptor(g.streamInterceptors...)
	g.options = append(g.options, interceptorOption, streamOpt)

	g.server = grpc.NewServer(g.options...)
}

func (g *grpcServer) Run() {
	if g.server == nil {
		g.server = grpc.NewServer(g.options...)
	}

	if !isPortAvailable(g.config.Port) {
		g.logger.Fatalf("gRPC port %d is blocked or unreachable", g.config.Port)
		return
	}

	addr := ":" + strconv.Itoa(g.config.Port)

	g.logger.Infof("starting gRPC server at %s", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		g.logger.Fatalf("error in starting gRPC server at %s: %s", addr, err)
		return
	}

	g.logger.Infof("gRPC server started successfully on %s", addr)

	if err := g.server.Serve(listener); err != nil {
		g.logger.Fatalf("error in serving gRPC server at %s: %s", addr, err)
		return
	}
}

func (g *grpcServer) Shutdown(ctx context.Context) error {
	if g.server == nil {
		return nil
	}

	ch := make(chan struct{}, 1)

	go func() {
		g.server.GracefulStop()
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		g.server.Stop()
	case <-ch:
	}

	return nil
}

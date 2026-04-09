package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mrbagir/qcash-appcore/examples/grpc-server/pb"
	appcore "github.com/mrbagir/qcash-appcore/pkg/app"
)

type Config struct {
	SERVICE_CLIENT_ADDRESS string `env:"GRPC_CLIENT" envDefault:"localhost:9001"`
}

type usecase struct {
	config Config

	pb.UnimplementedHelloServer
}

func (usecase) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

func main() {
	app := appcore.New()

	// Load .env
	var config Config
	app.ParseConfig(&config)

	usecase := &usecase{
		config: config,
	}

	// gRPC server
	pb.RegisterHelloServer(app, usecase)

	// HTTP server
	app.Handle("POST /api/sayhello", appcore.HandlerRPC(usecase.SayHello))

	app.Run()
}

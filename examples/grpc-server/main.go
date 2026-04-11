package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mrbagir/qcash-appcore/examples/grpc-server/pb"
	appcore "github.com/mrbagir/qcash-appcore/pkg/app"
)

type usecase struct {
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

	usecase := &usecase{}

	// gRPC server
	pb.RegisterHelloServer(app, usecase)

	// HTTP server
	app.Handle("POST /api/sayhello", appcore.HandlerRPC(usecase.SayHello))

	app.Run()
}

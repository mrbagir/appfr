package main

import (
	"context"

	"github.com/mrbagir/qcash-appcore/examples/grpc/pb"
	appcore "github.com/mrbagir/qcash-appcore/pkg/app"
)

type usecase struct {
	pb.UnimplementedHelloServer
}

func (usecase) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Client: Hello " + in.Name}, nil
}

func main() {
	app := appcore.New()

	pb.RegisterHelloServer(app, &usecase{})

	app.Run()
}

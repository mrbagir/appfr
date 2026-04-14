package main

import (
	"context"

	appcore "github.com/mrbagir/appfr"
	"github.com/mrbagir/appfr/examples/grpc-client/client/pb"
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

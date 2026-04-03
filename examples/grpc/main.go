package grpc

import (
	"context"

	"github.com/mrbagir/qcash-appcore/examples/grpc/pb"
	"github.com/mrbagir/qcash-appcore/pkg/app"
)

type usecase struct {
	pb.UnimplementedHelloServer
}

func (usecase) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error)

func main() {
	app := app.New()
	pb.RegisterHelloServer(app, &usecase{})
	app.Run()
}

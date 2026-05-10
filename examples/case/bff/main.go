package main

import (
	"github.com/mrbagir/appfr"
	"github.com/mrbagir/appfr/client"
	"github.com/mrbagir/appfr/examples/grpc-server/pb"
)

func main() {
	app := appfr.New()

	// Connect to the gRPC backend service
	helloClient := client.NewGRPCClient(app, ":9090", pb.NewHelloClient)

	// Expose HTTP endpoint for the frontend
	app.Handle("POST /api/sayhello", appfr.HandlerRPCOpts(helloClient.SayHello))

	app.Run()
}

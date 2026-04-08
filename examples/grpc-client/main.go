package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mrbagir/qcash-appcore/examples/grpc-client/client/pb"
	appcore "github.com/mrbagir/qcash-appcore/pkg/app"
	"github.com/mrbagir/qcash-appcore/pkg/client"
	"google.golang.org/grpc"
)

func main() {
	// Run client server
	cmd := exec.Command("go", "run", "./client")
	cmd.Env = append(os.Environ(), "GRPC_PORT=9001")
	_ = cmd.Start()

	app := appcore.New()

	// Connect to client server and call SayHello
	halloClient := client.NewGRPCClient(app, ":9001", pb.NewHelloClient)
	callClient(halloClient)

	app.Run()
}

func callClient(halloClient pb.HelloClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := halloClient.SayHello(ctx, &pb.HelloRequest{Name: "World"}, grpc.WaitForReady(true))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(res.Message)
}

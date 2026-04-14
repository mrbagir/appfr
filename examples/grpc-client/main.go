package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	appcore "github.com/mrbagir/appfr"
	"github.com/mrbagir/appfr/client"
	"github.com/mrbagir/appfr/examples/grpc-client/client/pb"
	"google.golang.org/grpc"
)

func main() {
	// Run client server
	cmd := runClientServer()
	defer cmd.Process.Kill()

	app := appcore.New()

	// Connect to client server and call SayHello
	helloClient := client.NewGRPCClient(app, ":9010", pb.NewHelloClient)
	callClient(helloClient)

	app.Run()
}

func callClient(helloClient pb.HelloClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := helloClient.SayHello(ctx, &pb.HelloRequest{Name: "World"}, grpc.WaitForReady(true))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(res.Message)
}

func runClientServer() *exec.Cmd {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "./client"
	_ = cmd.Start()
	return cmd
}

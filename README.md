# appfr

Framework Go sederhana yang dev-friendly dan tidak mengikat. Bangun HTTP server, gRPC server, gRPC client, dan cron jobs dengan setup minimal tanpa boilerplate berlebihan.

## Install

```bash
go get github.com/mrbagir/appfr
```

## Fitur

- **HTTP Server** — routing sederhana dengan gorilla/mux, support HTTPS
- **gRPC Server** — built-in recovery & observability interceptor
- **gRPC Client** — generic client dengan auto connection management
- **Cron Jobs** — scheduling dengan format cron 5/6 field
- **Structured Logging** — JSON & pretty print, multi-level
- **Graceful Shutdown** — otomatis handle signal interrupt
- **Config dari ENV** — parse `.env` file & environment variable

## Quick Start

### HTTP Server

```go
import "github.com/mrbagir/appfr"

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"message": "Hello!"})
}

func main() {
    app := appfr.New()
    app.Handle("POST /api/hello", HelloHandler)
    app.Run()
}
```

### gRPC Server

Contoh proto sederhana:

```proto
syntax = "proto3";

package hello;

option go_package = "./pb";

service Hello {
    rpc SayHello (HelloRequest) returns (HelloResponse);
}

message HelloRequest {
    string name = 1;
}

message HelloResponse {
    string message = 1;
}
```

```bash
protoc --go_out=. --go-grpc_out=. ./{path}/*.proto
```

Implementasi server:

```go
import "github.com/mrbagir/appfr"

type usecase struct{ pb.UnimplementedHelloServer }

func (u *usecase) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func main() {
    app := appfr.New()
    pb.RegisterHelloServer(app, &usecase{})
    app.Run()
}
```

gRPC + HTTP gateway sekaligus:

```go
import "github.com/mrbagir/appfr"

type usecase struct{ pb.UnimplementedHelloServer }

func (u *usecase) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func main() {
    app := appfr.New()
    uc := &usecase{}
    pb.RegisterHelloServer(app, uc)
    app.Handle("POST /api/hello", appfr.HandlerRPC(uc.SayHello))
    app.Run()
}
```

### gRPC Client

```go
import (
    "github.com/mrbagir/appfr"
    "github.com/mrbagir/appfr/client"
)

func main() {
    app := appfr.New()
    helloClient := client.NewGRPCClient(app, ":9010", pb.NewHelloClient)
    res, _ := helloClient.SayHello(ctx, &pb.HelloRequest{Name: "World"})
    fmt.Println(res.Message)
    app.Run()
}
```

### Cron Jobs

```go
import "github.com/mrbagir/appfr"

func main() {
    app := appfr.New()
    app.AddCronJob("* * * * * *", "every-second", func() {
        fmt.Println("tick")
    })
    app.Run()
}
```

> Contoh lengkap ada di folder [examples](examples/)

## Konfigurasi

Konfigurasi melalui environment variable atau file `.env`:

| Variable | Default | Deskripsi |
|---|---|---|
| `APP_ENV` | `development` | Environment aplikasi |
| `LOGGER_LEVEL` | `INFO` | Level log: DEBUG, INFO, NOTICE, WARN, ERROR, FATAL |
| `SERVER_SHUTDOWN_TIMEOUT` | `30s` | Timeout graceful shutdown |
| `HTTP_PORT` | `8080` | Port HTTP server |
| `HTTP_CERT_FILE` | - | Path certificate file (untuk HTTPS) |
| `HTTP_KEY_FILE` | - | Path key file (untuk HTTPS) |
| `GRPC_PORT` | `9000` | Port gRPC server |

Parse custom config:

```go
type MyConfig struct {
	DBHost string `env:"DB_HOST" envDefault:"localhost"`
}

app := appfr.New()

var cfg MyConfig
app.ParseConfig(&cfg)
```

## License

MIT

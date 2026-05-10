<div align="center">

# AppFr

**Framework Go sederhana yang dev-friendly dan tidak mengikat.**  
Bangun HTTP server, gRPC server, gRPC client, dan cron jobs — tanpa boilerplate berlebihan.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Go Reference](https://pkg.go.dev/badge/github.com/mrbagir/appfr.svg)](https://pkg.go.dev/github.com/mrbagir/appfr)
[![License](https://img.shields.io/badge/License-MIT-534AB7?style=flat-square)](LICENSE)

```bash
go get github.com/mrbagir/appfr@v1.1.0
```

</div>

---

## Kenapa AppFr?

AppFr menghilangkan hal-hal yang selalu ditulis berulang: setup graceful shutdown, wiring gRPC + HTTP gateway, parsing env config, dan scheduling cron — tanpa mendikte bagaimana kamu harus struktur aplikasimu. Kamu tetap pegang kendali penuh atas kode aplikasimu.

- **Tidak opinionated** — tidak ada magic, tidak ada convention yang dipaksakan
- **Minimal dependency** — hanya library battle-tested: `gorilla/mux`, `grpc`, `robfig/cron`
- **Satu entry point** — `app.Run()` menjalankan semua server sekaligus, blocking, dengan graceful shutdown otomatis saat menerima `SIGINT`/`SIGTERM`

---

## Fitur

- **HTTP Server** — routing dengan `gorilla/mux`, support HTTPS via cert file, compatible dengan `http.HandlerFunc` standar
- **gRPC Server** — panic recovery & logging interceptor bawaan via `go-grpc-middleware`, tanpa konfigurasi tambahan
- **gRPC + HTTP Gateway** — jalankan gRPC dan REST dari satu binary dengan `appfr.HandlerRPC`
- **gRPC Client** — generic client dengan auto connection management dan lazy initialization
- **Cron Jobs** — scheduling dengan format cron 5 atau 6 field (dengan detik)
- **Structured Logging** — JSON & pretty print, 6 level (DEBUG → FATAL), ganti level saat runtime
- **Graceful Shutdown** — handle `SIGINT`/`SIGTERM` otomatis dengan configurable timeout
- **Config dari ENV** — parse `.env` file & environment variable via struct tag `env:"..."`

---

## Quick Start

### HTTP Server

```go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/mrbagir/appfr"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"message": "Hello!"})
}

func main() {
    app := appfr.New()
    app.Handle("POST /api/hello", HelloHandler)
    app.Run() // blocking, graceful shutdown otomatis
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

Generate kode Go dari proto:

```bash
protoc --go_out=. --go-grpc_out=. ./{path}/*.proto
```

Implementasi server:

```go
package main

import (
    "context"

    "github.com/mrbagir/appfr"
    pb "github.com/yourorg/yourproto"
)

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

### gRPC + HTTP Gateway (dalam satu binary)

```go
func main() {
    app := appfr.New()

    uc := &usecase{} // usecase sudah embed pb.UnimplementedHelloServer
    pb.RegisterHelloServer(app, uc)

    // expose endpoint gRPC sebagai REST handler
    app.Handle("POST /api/hello", appfr.HandlerRPC(uc.SayHello))

    app.Run()
}
```

### gRPC Client

```go
package main

import (
    "context"
    "fmt"

    "github.com/mrbagir/appfr"
    "github.com/mrbagir/appfr/client"
    pb "github.com/yourorg/yourproto"
)

func main() {
    app := appfr.New()

    helloClient := client.NewGRPCClient(app, ":9010", pb.NewHelloClient)

    res, _ := helloClient.SayHello(context.Background(), &pb.HelloRequest{Name: "World"})
    fmt.Println(res.Message) // Hello World

    app.Run()
}
```

### Cron Jobs

```go
func main() {
    app := appfr.New()

    // format 6 field (dengan detik): detik menit jam hari bulan weekday
    app.AddCronJob("* * * * * *", "every-second", func() {
        fmt.Println("tick")
    })

    // format 5 field standar
    app.AddCronJob("0 9 * * 1-5", "weekday-morning", func() {
        fmt.Println("good morning!")
    })

    app.Run()
}
```

> Contoh lengkap ada di folder [`examples/`](./examples)

---

## Konfigurasi

Konfigurasi dibaca otomatis dari file `.env` (jika ada) dan environment variable. Tidak perlu setup manual.

| Variable | Default | Deskripsi |
|---|---|---|
| `APP_ENV` | `development` | Environment aplikasi |
| `LOGGER_LEVEL` | `INFO` | Level log: `DEBUG`, `INFO`, `NOTICE`, `WARN`, `ERROR`, `FATAL` |
| `SERVER_SHUTDOWN_TIMEOUT` | `30s` | Timeout graceful shutdown |
| `HTTP_PORT` | `8080` | Port HTTP server |
| `HTTP_CERT_FILE` | — | Path certificate file (untuk HTTPS) |
| `HTTP_KEY_FILE` | — | Path key file (untuk HTTPS) |
| `GRPC_PORT` | `9000` | Port gRPC server |

### Custom config

Gunakan `app.ParseConfig` untuk parse config struct milikmu sendiri. AppFr menggunakan [`caarlos0/env`](https://github.com/caarlos0/env) di balik layar.

```go
type MyConfig struct {
    DBHost string `env:"DB_HOST" envDefault:"localhost"`
    DBPort int    `env:"DB_PORT" envDefault:"5432"`
}

func main() {
    app := appfr.New()

    var cfg MyConfig
    app.ParseConfig(&cfg)

    fmt.Println(cfg.DBHost) // nilai dari env DB_HOST atau "localhost"

    app.Run()
}
```

### Mengakses logger

```go
app := appfr.New()
log := app.Logger()

log.Infof("server starting on port %d", 8080)
log.Debugf("config loaded: %+v", cfg)
```

---

## Cara `app.Run()` bekerja

`app.Run()` bersifat **blocking** — ia menjalankan semua server yang terdaftar (HTTP, gRPC, cron) secara concurrent lalu menunggu hingga semua selesai.

Saat menerima sinyal `SIGINT` atau `SIGTERM` (misalnya dari `Ctrl+C` atau `kill`), AppFr secara otomatis:

1. Menghentikan penerimaan request baru
2. Menunggu request yang sedang berjalan selesai (hingga batas `SERVER_SHUTDOWN_TIMEOUT`)
3. Menutup semua gRPC client connection
4. Menghentikan cron scheduler

Tidak perlu menulis shutdown logic sendiri.

---

## Lisensi

MIT — lihat [LICENSE](./LICENSE)

---

## Kontribusi

Pull request dan issue sangat disambut. Untuk perubahan besar, buka issue dulu agar kita bisa diskusi arah pengembangannya.
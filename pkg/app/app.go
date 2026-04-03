package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mrbagir/qcash-appcore/pkg/config"
)

type App struct {
	grpcServer *grpcServer

	grpcRegistered bool
	httpRegistered bool

	config config.Config
}

func New() *App {
	app := &App{}
	app.grpcServer = newGRPCServer(app.config)

	return app
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.startAllServers(ctx)
}

func (a *App) startAllServers(ctx context.Context) {
	wg := sync.WaitGroup{}

	// a.startMetricsServer(&wg)
	// a.startHTTPServer(&wg)
	a.startGRPCServer(&wg)

	wg.Wait()
}

func (a *App) startGRPCServer(wg *sync.WaitGroup) {
	if a.grpcRegistered {
		wg.Add(1)

	}
}

package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mrbagir/qcash-appcore/pkg/logging"
)

type config struct {
	AppEnv          string        `env:"APP_ENV" envDefault:"development"`
	LoggerLevel     string        `env:"LOGGER_LEVEL" envDefault:"INFO"`
	ShutdownTimeout time.Duration `env:"TIMEOUT" envDefault:"30s"`
	HttpConfig      httpConfig
	GrpcConfig      grpcConfig
}

type App struct {
	grpcServer *grpcServer
	httpServer *httpServer

	grpcRegistered bool
	httpRegistered bool

	grpcClients *grpcClients

	config config
	logger logging.Logger
}

func New() *App {
	app := &App{}
	app.logger = logging.NewLogger(logging.INFO)
	app.readConfig()
	app.ParseConfig(&app.config)
	app.logger.ChangeLevel(logging.GetLevelFromString(app.config.LoggerLevel))

	app.httpServer = newHTTPServer(app.logger, app.config)
	app.grpcServer = newGRPCServer(app.logger, app.config)

	return app
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.startShutdownHandler(ctx)
	a.startAllServers(ctx)
}

func (a *App) startShutdownHandler(ctx context.Context) {
	go func() {
		<-ctx.Done()

		shutdownCtx, done := context.WithTimeout(context.WithoutCancel(ctx), a.config.ShutdownTimeout)
		defer done()

		a.logger.Infof("Shutting down server with a timeout of %v", a.config.ShutdownTimeout)

		shutdownErr := a.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			a.logger.Debugf("Server shutdown failed: %v", shutdownErr)
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) error {
	var err error
	if a.httpServer != nil {
		err = errors.Join(err, a.httpServer.Shutdown(ctx))
	}

	if a.grpcServer != nil {
		err = errors.Join(err, a.grpcServer.Shutdown(ctx))
	}

	if a.grpcClients != nil {
		a.grpcClients.shutdown()
	}

	return err
}

func (a *App) startAllServers(_ context.Context) {
	wg := sync.WaitGroup{}

	// a.startMetricsServer(&wg)
	a.startHTTPServer(&wg)
	a.startGRPCServer(&wg)

	wg.Wait()
}

func (a *App) startHTTPServer(wg *sync.WaitGroup) {
	if a.httpRegistered {
		wg.Add(1)
		go func(s *httpServer) {
			defer wg.Done()
			s.Run()
		}(a.httpServer)
	}
}

func (a *App) startGRPCServer(wg *sync.WaitGroup) {
	if a.grpcRegistered {
		wg.Add(1)
		go func(s *grpcServer) {
			defer wg.Done()
			s.Run()
		}(a.grpcServer)
	}
}

func (a *App) Logger() logging.Logger {
	return a.logger
}

func isPortAvailable(port int) bool {
	dialer := net.Dialer{Timeout: checkPortTimeout}

	conn, err := dialer.DialContext(context.Background(), "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}

	conn.Close()

	return false
}

package appfr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"

	"github.com/mrbagir/appfr/http/middleware"
	"github.com/mrbagir/appfr/logging"
)

type httpConfig struct {
	Port     int    `env:"HTTP_PORT" envDefault:"8080"`
	CertFile string `env:"HTTP_CERT_FILE"`
	KeyFile  string `env:"HTTP_KEY_FILE"`
}

type httpServer struct {
	server *http.Server
	router *mux.Router
	config httpConfig
	logger logging.Logger
}

var (
	errInvalidCertificateFile = errors.New("invalid certificate file")
	errInvalidKeyFile         = errors.New("invalid key file")
)

func (a *App) Handle(path string, handler http.HandlerFunc) {
	parts := strings.SplitN(path, " ", 2)
	switch len(parts) {
	case 1:
		a.httpServer.router.Handle(path, handler)
	case 2:
		a.httpServer.router.Handle(parts[1], handler).Methods(strings.Split(parts[0], ",")...)
	}
	a.httpRegistered = true
}

var decoder = schema.NewDecoder()

func HandlerRPC[IN, OUT any](fn func(context.Context, *IN) (*OUT, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in IN
		if err := decoder.Decode(&in, r.URL.Query()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request body: %v"}`, err)
			return
		}

		out, err := fn(r.Context(), &in)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
				fmt.Fprintf(w, `{"error":"%s"}`, s.Message())
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	}
}

func newHTTPServer(logger logging.Logger, cfg config) *httpServer {
	r := mux.NewRouter()

	r.Use(
		middleware.Logging(logger),
	)

	return &httpServer{
		router: r,
		config: cfg.HttpConfig,
		logger: logger,
	}
}

func (h *httpServer) Run() {
	if h.server != nil {
		h.logger.Warnf("Server already running on port: %d", h.config.Port)
		return
	}

	h.logger.Infof("starting HTTP server at port %d", h.config.Port)

	h.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", h.config.Port),
		Handler:           h.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// If both certFile and keyFile are provided, validate and run HTTPS server
	if h.config.CertFile != "" && h.config.KeyFile != "" {
		if err := validateCertificateAndKeyFiles(h.config.CertFile, h.config.KeyFile); err != nil {
			h.logger.Error(err)
			return
		}

		// Start HTTPS server with TLS
		if err := h.server.ListenAndServeTLS(h.config.CertFile, h.config.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.logger.Errorf("error while listening to https server, err: %v", err)
		}

		return
	}

	// If no certFile/keyFile is provided, run the HTTP server
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		h.logger.Errorf("error while listening to http server, err: %v", err)
	}
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Close()
}

func validateCertificateAndKeyFiles(certificateFile, keyFile string) error {
	if _, err := os.Stat(certificateFile); os.IsNotExist(err) {
		return fmt.Errorf("%w : %v", errInvalidCertificateFile, certificateFile)
	}

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return fmt.Errorf("%w : %v", errInvalidKeyFile, keyFile)
	}

	return nil
}

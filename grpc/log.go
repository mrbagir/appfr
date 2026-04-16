package grpc

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mrbagir/appfr/internal/trace"
)

const (
	statusCodeWidth   = 3
	responseTimeWidth = 11
)

type gRPCLog struct {
	ID           string `json:"id"`
	StartTime    string `json:"startTime"`
	ResponseTime int64  `json:"responseTime"`
	Method       string `json:"method"`
	StatusCode   int32  `json:"statusCode"`
	StreamType   string `json:"streamType,omitempty"`
}

func (l *gRPCLog) PrettyPrint(writer io.Writer) {
	streamInfo := ""
	if l.StreamType != "" {
		streamInfo = fmt.Sprintf(" [%s]", l.StreamType)
	}

	fmt.Fprintf(writer, "\u001B[38;5;8m%s \u001B[38;5;%dm%-*d"+
		"\u001B[0m %*d\u001B[38;5;8mµs\u001B[0m %s%s %s\n",
		l.ID, colorForGRPCCode(l.StatusCode),
		statusCodeWidth, l.StatusCode,
		responseTimeWidth, l.ResponseTime,
		"GRPC", streamInfo, l.Method)
}

func colorForGRPCCode(s int32) int {
	const (
		blue = 34
		red  = 202
	)

	if s == 0 {
		return blue
	}

	return red
}

type Logger interface {
	Info(args ...any)
	Errorf(string, ...any)
	Debug(...any)
	Fatalf(string, ...any)
}

func ObservabilityInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)
		if err != nil && isServerError(err) {
			logger.Errorf("error while handling gRPC request to method %q: %q", info.FullMethod, err)
		}

		logRPC(ctx, logger, start, err, info.FullMethod)

		return resp, err
	}
}

func logRPC(ctx context.Context, logger Logger, start time.Time, err error, method string) {
	duration := time.Since(start)

	logEntry := gRPCLog{
		ID:           trace.GetProcessIDFromCtx(ctx),
		StartTime:    start.Format("2006-01-02T15:04:05.999999999-07:00"),
		ResponseTime: duration.Microseconds(),
		Method:       method,
	}

	if err != nil {
		statusErr, _ := status.FromError(err)
		logEntry.StatusCode = int32(statusErr.Code())
	} else {
		logEntry.StatusCode = int32(codes.OK)
	}

	logger.Info(&logEntry)
}

// isServerError returns true if the gRPC error represents a server-side error.
// Client errors like ResourceExhausted, InvalidArgument, NotFound, etc. are not
// considered server errors and should not be logged at ERROR level.
func isServerError(err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return true
	}

	switch s.Code() {
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.Unauthenticated, codes.ResourceExhausted,
		codes.FailedPrecondition, codes.OutOfRange, codes.Canceled:
		return false
	case codes.OK, codes.Unknown, codes.DeadlineExceeded, codes.Aborted,
		codes.Unimplemented, codes.Internal, codes.Unavailable, codes.DataLoss:
		return true
	default:
		return true
	}
}

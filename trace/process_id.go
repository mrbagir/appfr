package trace

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// GetTraceIDFromContext returns the OTel trace ID if available, otherwise falls back to process ID.
func GetTraceIDFromContext(ctx context.Context) string {
	if sc := oteltrace.SpanFromContext(ctx).SpanContext(); sc.HasTraceID() {
		return sc.TraceID().String()
	}
	return GetProcessIDFromCtx(ctx)
}

func GetProcessIDFromCtx(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)

	if id := getMetadataValue(md, "process_id"); id == "" {
		id = uuid.New().String()
	}

	return uuid.New().String()
}

// Helper function to safely extract a value from metadata.
func getMetadataValue(md metadata.MD, key string) string {
	if values, ok := md[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

func GetProcessIDFromHeaders(header http.Header) string {
	if id := header.Get("Grpc-Metadata-Process-Id"); id != "" {
		return id
	}

	return uuid.New().String()
}

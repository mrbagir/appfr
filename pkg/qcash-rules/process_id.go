package qcashrules

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func GetProcessID(ctx context.Context) string {
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

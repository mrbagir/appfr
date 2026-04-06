package grpc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

var FromError = status.FromError
var HTTPStatusFromCode = runtime.HTTPStatusFromCode

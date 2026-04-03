package runner

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const metadataTraceID = "x-trace-id"

func incomingTraceID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	v := md.Get(metadataTraceID)
	if len(v) == 0 {
		return ""
	}

	return strings.TrimSpace(v[0])
}

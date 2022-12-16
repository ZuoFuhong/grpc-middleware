package tracing

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// UnaryServerInterceptor returns a new unary server interceptor for Tracing.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if metadata.ValueFromIncomingContext(ctx, TraceId) == nil {
			// 生成新的 traceId
			traceId := uuid.New().String()
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				md.Set(TraceId, traceId)
				ctx = metadata.NewIncomingContext(ctx, md)
			}
		}
		start := time.Now()
		resp, err := handler(ctx, req)
		// Report trace log
		go dcReport(ctx, info.FullMethod, req, start, err)
		return resp, err
	}
}

package tracing

import (
	"context"
	"github.com/ZuoFuhong/grpc-middleware/dc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// UnaryClientInterceptor returns a new unary client interceptor for Tracing.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		// 通过 context 透传到下游
		meta, _ := metadata.FromIncomingContext(ctx)
		err := invoker(metadata.NewOutgoingContext(ctx, meta), method, req, reply, cc, opts...)
		// Report trace log
		dc.Report(ctx, method, req, start, err)
		return err
	}
}

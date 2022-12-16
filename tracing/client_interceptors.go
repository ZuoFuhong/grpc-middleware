package tracing

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// UnaryClientInterceptor returns a new unary client interceptor for Tracing.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		// Report trace log
		go dcReport(ctx, method, req, start, err)
		return err
	}
}

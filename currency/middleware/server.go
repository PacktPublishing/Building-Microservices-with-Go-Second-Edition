package middleware

import (
	"context"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewUnaryServerRequestLogger(logger hclog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		switch req.(type) {
		case *currency.RateRequest:
			rr := req.(*currency.RateRequest)
			logger.Info("New Request", "method", info.FullMethod, "base", rr.GetBase(), "destination", rr.GetDestination(), "metadata", md)
		}

		resp, err := handler(ctx, req)
		return resp, err
	}
}

func NewStreamingServerRequestLogger(logger hclog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger.Info("Method called", "name", info.FullMethod, "client_stream", info.IsClientStream, "server_stream", info.IsServerStream)

		return handler(srv, ss)
	}
}

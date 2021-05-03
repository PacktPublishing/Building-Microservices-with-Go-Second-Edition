package middleware

import (
	"context"
	"time"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewUnaryClientRequestLogger(logger hclog.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, _ := metadata.FromIncomingContext(ctx)

		switch req.(type) {
		case *currency.RateRequest:
			rReq := req.(*currency.RateRequest)
			st := time.Now()
			logger.Info("New unary request", "method", method, "base", rReq.GetBase(), "destination", rReq.GetDestination(), "metadata", md)

			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return err
			}

			rResp := reply.(*currency.RateResponse)
			logger.Info("Got unary reply", "method", method, "rate", rResp.Rate, "duration", time.Now().Sub(st), "metadata", md)
		}

		return nil
	}
}

func NewStreamingClientRequestLogger(logger hclog.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		logger.Info("New stream request", "method", method)

		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}

		return clientStream, nil
	}
}

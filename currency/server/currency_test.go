package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestClient(t *testing.T) {
	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := currency.NewCurrencyClient(conn)
	md := metadata.New(
		map[string]string{
			"test": "123",
		},
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	rr, err := c.GetRate(ctx, &currency.RateRequest{Base: "EUR", Destination: "GBP"})
	if err != nil {
		panic(err)
	}

	fmt.Println(rr)
}

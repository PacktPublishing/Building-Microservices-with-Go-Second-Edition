package server

import (
	"context"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/metadata"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	log hclog.Logger
}

// NewCurrency creates a new Currency server
func NewCurrency(l hclog.Logger) *Currency {
	c := &Currency{l}

	return c
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination(), "metadata", md)

	return &currency.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: 1.25}, nil
}

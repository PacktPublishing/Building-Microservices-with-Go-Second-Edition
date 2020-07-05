package server

import (
	"context"
	"fmt"
	"io"

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

func (c *Currency) SubscribeRates(svr currency.Currency_SubscribeRatesServer) error {

	// handle client requests
	receiveError := c.handleClientMessages(svr)

	// send messages to the client
	sendError := c.handleServerMessages(svr)

	// block until we get an error from the receive or the send loop
	select {
	case err := <-receiveError:
		return err
	case err := <-sendError:
		return err
	}
}

func (c *Currency) handleClientMessages(svr currency.Currency_SubscribeRatesServer) chan error {
	receiveError := make(chan error)

	go func() {
		for {
			rr, err := svr.Recv()
			// client has closed the connection
			if err == io.EOF {
				c.log.Error("Client disconnected, closing connection")
				receiveError <- nil
			}

			// generic error send a reply to the client
			if err != nil {
				c.log.Error("Error recieved when reading from client stream", "error", err)
				receiveError <- fmt.Errorf("Unable to read message", "error", err)
			}

			// process the request
			c.log.Info("Received message from client", "base", rr.GetBase(), "dest", rr.GetDestination())
		}
	}()

	return receiveError
}

func (c *Currency) handleServerMessages(svr currency.Currency_SubscribeRatesServer) chan error {
	sendError := make(chan error)

	return sendError
}

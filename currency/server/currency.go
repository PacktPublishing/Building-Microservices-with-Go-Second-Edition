package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	// if the base and the destination currency is the same return an error
	if rr.Base == rr.Destination {
		// create a new gRPC status
		statusErr := status.New(codes.InvalidArgument, "Base currency can not be the same as the destination currency")

		// add the request to status payload
		var err error
		statusErr, err = statusErr.WithDetails(rr)
		if err != nil {
			// Unable to add request to status payload, return an error, this should never happen
			c.log.Error("Unable to add request to status details", "request", rr, "error", err)
			return nil, err
		}

		return nil, statusErr.Err()
	}

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

	go func() {
		for {
			// send a message to the client
			err := svr.Send(&currency.RateResponse{Base: currency.Currencies_EUR, Destination: currency.Currencies_USD, Rate: 1.25})
			if err != nil {
				c.log.Error("Unable to send message to the client", "error", err)
				sendError <- err
			}

			// sleep for 5s before retrying
			time.Sleep(5 * time.Second)
		}
	}()

	return sendError
}

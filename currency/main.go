package main

import (
	"fmt"
	"net"
	"os"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

func main() {
	log := hclog.Default()

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// create an instance of the Currency server
	c := server.NewCurrency(log)

	// register the currency server
	currency.RegisterCurrencyServer(gs, c)

	// create a TCP socket for inbound server connections
	serverAddress := fmt.Sprintf(":%d", 9092)
	l, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	// listen for requests
	log.Info("Starting Currency server", "address", serverAddress)
	gs.Serve(l)
}

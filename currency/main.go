package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// create an instance of the Currency server
	c := server.NewCurrency(log)

	// register the currency server
	currency.RegisterCurrencyServer(gs, c)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	serverAddress := fmt.Sprintf(":%d", 9092)
	l, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	// listen for requests
	go func() {
		log.Info("Starting Currency server", "address", serverAddress)
		gs.Serve(l)
	}()

	log.Info("Press Ctrl-C to stop service")
	stopChan := make(chan os.Signal, 1)
	// register for sig kill and sig interupts
	signal.Notify(stopChan, os.Interrupt)
	signal.Notify(stopChan, os.Kill)

	// Block until a signal is received.
	sig := <-stopChan
	log.Info("Caught signal, waiting for connections to exit before shutting down server", "signal", sig)
	gs.GracefulStop()
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/4_handlers/handlers"
)

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	// create the handlers
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	// create a new serve mux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	// Listen for connections on all ip addresses (0.0.0.0)
	// port 9090
	log.Println("Starting Server")
	err := http.ListenAndServe(":9090", sm)
	log.Fatal(err)
}

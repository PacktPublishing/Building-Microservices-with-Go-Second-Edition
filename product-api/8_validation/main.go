package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/7_Gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()

	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	// create the handlers
	ph := handlers.NewProducts(l)

	// create a new Gorilla mux router and register the handlers
	sm := mux.NewRouter()

	getR := sm.Methods("GET").Subrouter()
	getR.HandleFunc("/products", ph.GET)

	putR := sm.Methods("PUT").Subrouter()
	putR.HandleFunc("/products/{id:[0-9]+}", ph.PUT)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods("POST").Subrouter()
	postR.HandleFunc("/products", ph.POST)
	postR.Use(ph.MiddlewareValidateProduct)

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	//go func() {
	l.Println("[INFO] Starting server on port 9090")

	err := s.ListenAndServe()
	if err != nil {
		l.Printf("[ERROR] Error starting server: %s\n", err)
		os.Exit(1)
	}
	//}()

	// trap sigterm or interupt and gracefully shutdown the server
	l.Println("[INFO] Press Ctrl-C to stop service")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("[INFO] Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/9_docs/data"
	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/9_docs/handlers"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"

	"github.com/go-openapi/runtime/middleware"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
var allowedOrigins = env.String("ALLOW_ORIGIN", false, "http://localhost:3000", "Allowe origin for CORS requests")

func main() {

	env.Parse()

	l := log.New(os.Stdout, "products-api ", log.LstdFlags)
	v := data.NewValidation()

	// create the handlers
	ph := handlers.NewProducts(l, v)

	// create a new Gorilla mux router and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.ListAll)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)
	getR.Use(ph.MiddlewareContentType)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products", ph.Update)
	putR.Use(ph.MiddlewareContentType)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddlewareContentType)
	postR.Use(ph.MiddlewareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)
	deleteR.Use(ph.MiddlewareContentType)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getR.Handle("/docs", sh)

	// handler to return the swagger documentation
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// allow the local we
	ch := ghandlers.CORS(
		ghandlers.AllowedOrigins([]string{*allowedOrigins}),
	)(sm)

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch,                // set the default handler
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
	// always need to call cancel to avoid leaking context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}

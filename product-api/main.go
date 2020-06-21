package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	protos "github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/data"
	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/handlers"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"google.golang.org/grpc"

	"github.com/go-openapi/runtime/middleware"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
var currencyAddress = env.String("CURRENCY_ADDRESS", false, "localhost:9092", "Address for the currency server")
var allowedOrigins = env.String("ALLOW_ORIGIN", false, "http://localhost:3000", "Allowed origin for CORS requests")

func main() {

	env.Parse()

	l := hclog.Default().Named("product-api")

	v := data.NewValidation()

	cc, closeFunc := createCurrencyClient(l)
	defer closeFunc()

	db := data.NewProductsDB(cc, l)

	// create the handlers
	ch := createHandlers(v, db, l)

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,                                     // configure the bind address
		Handler:      ch,                                               // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server", "address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	l.Info("Press Ctrl-C to stop service")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Caught signal:", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	// always need to call cancel to avoid leaking context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}

func createCurrencyClient(l hclog.Logger) (protos.CurrencyClient, func()) {
	// create the currency service client
	conn, err := grpc.Dial(*currencyAddress, grpc.WithInsecure())
	if err != nil {
		l.Info("Unable to create client for currency service", "error", err)
		os.Exit(1)
	}

	return protos.NewCurrencyClient(conn), func() {
		defer conn.Close()
	}
}

func createHandlers(v *data.Validation, db *data.ProductsDB, l hclog.Logger) http.Handler {
	ph := handlers.NewProducts(db, v, l.Named("products-handler"))

	// create a new Gorilla mux router and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.ListAll)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)
	getR.HandleFunc("/products", ph.ListAll).Queries("currency", "{[A-Z]{3}}")
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle).Queries("currency", "{[A-Z]{3}}")
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

	return ch
}

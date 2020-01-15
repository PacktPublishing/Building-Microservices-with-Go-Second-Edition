package handlers

import (
	"log"
	"net/http"
	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/6_REST/data"
)

// Products handler for getting and updating products
type Products struct {
	l *log.Logger
}

// NewProducts returns a new products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP implements the http.Handler interface
func (p*Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	prods := data.GetProducts()
	
	err := prods.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Error serialzing products", http.StatusInternalServerError)
	}
}

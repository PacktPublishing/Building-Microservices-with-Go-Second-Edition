package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/9_docs/data"
	"github.com/gorilla/mux"
)

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

// Products handler for getting and updating products
type Products struct {
	l *log.Logger
	v *data.Validation
}

// NewProducts returns a new products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message error `json:"message"`
}

// getProductID returns the product ID from the URL
func getProductID(r *http.Request) (int, error) {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	return strconv.Atoi(vars["id"])
}

// GET handles GET requests and returns all current products
func (p *Products) GET(rw http.ResponseWriter, r *http.Request) {
	prods := data.GetProducts()

	err := data.ToJSON(prods, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
		http.Error(rw, "Error serialzing products", http.StatusInternalServerError)
	}
}

// PUT handles PUT requests to update products
func (p *Products) PUT(rw http.ResponseWriter, r *http.Request) {
	// fetch the id from the query string
	id, err := getProductID(r)
	if err != nil {
		p.l.Println("[ERROR] unable to find product id in URL", r.URL.Path, err)
		http.Error(rw, "Missing product id, url should be formatted /products/[id] for PUT requests", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	// override the product id
	prod.ID = id

	err = data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)
		http.Error(rw, "Product not found in database", http.StatusNotFound)
		return
	}

	p.l.Printf("[DEBUG] Updated product: %#v\n", prod)

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}

// POST handles post requests to add new products
func (p *Products) POST(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(prod)

	p.l.Printf("[DEBUG] Inserted product: %#v\n", prod)
}

// MiddlewareValidateProduct validates the product in the request and calls next if ok
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := data.FromJSON(prod, r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = p.v.Validate(prod)
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(rw, "Error validating product", http.StatusPreconditionFailed)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

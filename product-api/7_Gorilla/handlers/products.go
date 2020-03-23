package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/7_Gorilla/data"
	"github.com/gorilla/mux"
)

// unexported context key
type keyProduct struct{}

// Products handler for getting and updating products
type Products struct {
	l *log.Logger
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// NewProducts returns a new products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// getProductID returns the product ID from the URL
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}

// ListProducts handles HTTP GET requests for the products
func (p *Products) ListProducts(rw http.ResponseWriter, r *http.Request) {
	prods := data.GetProducts()

	err := data.ToJSON(prods, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
	}
}

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Println("[DEBUG] get record id", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}
}

// UpdateProduct handles PUT requests to update products
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// get the product from the context
	// middleware decodes this from the http.Request
	prod := r.Context().Value(keyProduct{}).(data.Product)

	err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: fmt.Sprintf("Product %d not found in database", prod.ID)})
		return
	}

	p.l.Printf("[DEBUG] Updated product: %#v\n", prod)

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}

// CreateProduct handles post requests to add new products
func (p *Products) CreateProduct(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(keyProduct{}).(data.Product)
	data.AddProduct(prod)

	p.l.Printf("[DEBUG] Inserted product: %#v\n", prod)
}

// MiddlewareValidateProduct validates the product in the request and calls next if ok
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := data.FromJSON(&prod, r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), keyProduct{}, prod)
		r = r.Clone(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

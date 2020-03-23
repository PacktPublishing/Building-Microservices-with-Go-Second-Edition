package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

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

// ErrInvalidProductPath error message
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// ServeHTTP implements the http.Handler interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, _ := getProductID(r)
		if id > 0 {
			p.getSingle(id, rw, r)
			return
		}

		p.getAll(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.put(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.post(rw, r)
		return
	}

	if r.Method == http.MethodDelete {
		p.delete(rw, r)
		return
	}

	// if we have not matched any of the HTTP methods handled return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProductID returns the product ID from the URL
func getProductID(r *http.Request) (int, error) {
	// Parse the product id from the URI
	re := regexp.MustCompile(`^\/products\/([0-9]+)$`)
	m := re.FindAllStringSubmatch(r.URL.Path, -1)

	// We should have one match which contains two groups
	// anything else is an invalid URI
	if len(m) != 1 || len(m[0]) != 2 {
		return -1, ErrInvalidProductPath
	}

	// Convert the id into an integer and return
	// the regex ensures that the second group is an integer
	return strconv.Atoi(m[0][1])
}

// getSingle handles HTTP GET requests for the products returning a single product
func (p *Products) getSingle(id int, rw http.ResponseWriter, r *http.Request) {
	prod, err := data.GetProductByID(id)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)
		http.Error(rw, "Product not found in database", http.StatusNotFound)
		return
	}

	err = data.ToJSON(rw, prod)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}
}

// getAll handles HTTP GET requests for the products returning all products
func (p *Products) getAll(rw http.ResponseWriter, r *http.Request) {

	prods := data.GetProducts()

	err := data.ToJSON(rw, prods)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}
}

func (p *Products) put(rw http.ResponseWriter, r *http.Request) {
	prod := data.Product{}

	// deserialize the body
	err := data.FromJSON(r.Body, &prod)
	if err != nil {
		p.l.Println("[ERROR] deserializing product", err)
		http.Error(rw, "Error reading product", http.StatusBadRequest)
		return
	}

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

func (p *Products) post(rw http.ResponseWriter, r *http.Request) {
	prod := data.Product{}

	err := data.FromJSON(r.Body, &prod)
	if err != nil {
		p.l.Println("[ERROR] deserializing product", err)
		http.Error(rw, "Error reading product", http.StatusBadRequest)
		return
	}

	data.AddProduct(prod)

	p.l.Printf("[DEBUG] Inserted product: %#v\n", prod)

	// return the product with the inserted id

}

func (p *Products) delete(rw http.ResponseWriter, r *http.Request) {
	// fetch the id from the query string
	id, err := getProductID(r)
	if err != nil {
		p.l.Println("[ERROR] unable to find product id in URL", r.URL.Path, err)
		http.Error(rw, "Missing product id, url should be formatted /products/[id] for PUT requests", http.StatusBadRequest)
		return
	}

	err = data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)
		http.Error(rw, "Product not found in database", http.StatusNotFound)
		return
	}

	p.l.Printf("[DEBUG] Deleted product: %#v\n", id)

	rw.WriteHeader(http.StatusNoContent)
}

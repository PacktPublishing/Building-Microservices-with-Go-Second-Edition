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
		p.get(rw, r)
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

	// if we have not matched any of the HTTP methods handled return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProductID returns the product ID from the URL
func getProductID(r *http.Request) (int, error) {
	// parse the product id from the url
	re := regexp.MustCompile(`\/products\/([0-9]+)`)
	m := re.FindAllStringSubmatch(r.URL.Path, -1)

	// if there is more than one match the URL is invalid
	if len(m) != 1 {
		return -1, ErrInvalidProductPath
	}

	// there should be two match groups the second contains the id
	if len(m[0]) != 2 {
		return -1, ErrInvalidProductPath
	}

	// convert the id into an integer and return
	return strconv.Atoi(m[0][1])
}

// get handles HTTP GET requests for the products
func (p *Products) get(rw http.ResponseWriter, r *http.Request) {
	prods := data.GetProducts()

	err := prods.ToJSON(rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
		http.Error(rw, "Error serialzing products", http.StatusInternalServerError)
	}
}

func (p *Products) put(rw http.ResponseWriter, r *http.Request) {
	prod := data.Product{}

	// fetch the id from the query string
	id, err := getProductID(r)
	if err != nil {
		p.l.Println("[ERROR] unable to find product id in URL", r.URL.Path, err)
		http.Error(rw, "Missing product id, url should be formatted /products/[id] for PUT requests", http.StatusBadRequest)
		return
	}

	// deserialize the body
	err = prod.FromJSON(r.Body)
	if err != nil {
		p.l.Println("[ERROR] deserializing product", err)
		http.Error(rw, "Error reading product", http.StatusBadRequest)
		return
	}

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

func (p *Products) post(rw http.ResponseWriter, r *http.Request) {
	prod := data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		p.l.Println("[ERROR] deserializing product", err)
		http.Error(rw, "Error reading product", http.StatusBadRequest)
		return
	}

	data.AddProduct(prod)

	p.l.Printf("[DEBUG] Inserted product: %#v\n", prod)

	// return the product with the inserted id

}

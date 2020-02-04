package handlers

import (
	"net/http"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/9_docs/data"
)

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	prods := data.GetProducts()

	err := data.ToJSON(prods, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
		http.Error(rw, "Error serialzing products", http.StatusInternalServerError)
	}
}

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	prods := data.GetProducts()

	err := data.ToJSON(prods, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
		http.Error(rw, "Error serialzing products", http.StatusInternalServerError)
	}
}

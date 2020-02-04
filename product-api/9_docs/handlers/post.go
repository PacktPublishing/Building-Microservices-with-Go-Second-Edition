package handlers

import (
	"net/http"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/9_docs/data"
)

// Create handles POST requests to add new products
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(prod)

	p.l.Printf("[DEBUG] Inserted product: %#v\n", prod)
}

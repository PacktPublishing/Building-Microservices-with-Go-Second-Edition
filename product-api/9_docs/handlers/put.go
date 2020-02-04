package handlers

import (
	"net/http"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/9_docs/data"
)

// Update handles PUT requests to update products
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
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

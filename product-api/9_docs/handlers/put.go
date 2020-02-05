package handlers

import (
	"net/http"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/9_docs/data"
)

// Update handles PUT requests to update products
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	// fetch the id from the query string
	id := getProductID(r)

	p.l.Println("[DEBUG] updating record id", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	// override the product id
	prod.ID = id

	err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	p.l.Printf("[DEBUG] Updated product: %#v\n", prod)

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}

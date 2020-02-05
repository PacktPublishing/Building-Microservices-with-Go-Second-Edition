package handlers

import (
	"net/http"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/9_docs/data"
)

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] get all records")

	prods := data.GetProducts()

	err := data.ToJSON(prods, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
		http.Error(rw, "Error serialzing products", http.StatusInternalServerError)
	}
}

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Println("[DEBUG] get record id", id)

	prod, err := data.GetProductByID(id)
	if err != nil {
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != data.ErrProductNotFound {
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
		http.Error(rw, "Error serialzing products", http.StatusInternalServerError)
	}
}

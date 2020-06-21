package handlers

import (
	"net/http"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/data"
)

// swagger:route GET /products products listAll
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	cur := r.URL.Query().Get("currency")
	p.l.Debug("Get all records", "curency", cur)

	rw.Header().Add("Content-Type", "application/json")

	prods, err := p.db.GetProducts(cur)
	if err != nil {
		p.l.Error("Unable to get products")

	}

	err = data.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serializing product", "error", err)
	}
}

// swagger:route GET /products/{id} products listSingle
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	cur := r.URL.Query().Get("currency")
	p.l.Debug("Get record", "id", id, "currency", cur)

	rw.Header().Add("Content-Type", "application/json")

	prod, err := p.db.GetProductByID(id, cur)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Product not found", "id", id)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("Unable to fetch product", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serialize product", "error", err)
	}
}

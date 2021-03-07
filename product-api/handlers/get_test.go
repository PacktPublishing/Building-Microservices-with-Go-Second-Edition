package handlers

import (
	"testing"

	protos "github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/data"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

func setupProductsHandler(t *testing.T) {
	l := hclog.Default().Named("product-api")
	v := data.NewValidation()
	conn, _ := grpc.Dial("localhost:9092", grpc.WithInsecure())
	cc := protos.NewCurrencyClient(conn)

	db, err := data.NewProductsDB(cc, l)
	ph := handlers.NewProducts(db, v, l.Named("products-handler"))

	t.Cleanup(func() {

	})
}

func TestReturnsAllProducts(t *testing.T) {

}

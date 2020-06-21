package main

import (
	"fmt"
	"os"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/client-sdk/client"
	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/client-sdk/client/products"
)

func main() {
	fmt.Println("Test App for Swagger generated client")

	t := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, t)

	p, err := c.Products.ListAll(products.NewListAllParams())
	if err != nil {
		fmt.Println("Unable to get list products, error:", err)
		os.Exit(1)
	}

	m := p.GetPayload()
	for _, m := range m {
		fmt.Printf("Product ID: %d Name: %s Price: %.2f SKU: %s\n", m.ID, *m.Name, *m.Price, *m.SKU)
	}
}

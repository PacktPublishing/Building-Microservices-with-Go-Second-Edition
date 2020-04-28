package data

import (
	"fmt"

	protos "github.com/nicholasjackson/building-microservices-youtube/currency/protos/currency"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// DB defines a product data model for CRUD based operations
type DB struct {
	cc protos.CurrencyClient
}

// NewDB creates a new data model
func NewDB(cc protos.CurrencyClient) *DB {
	return &DB{cc}
}

// GetProducts returns all products from the database
func (db *DB) GetProducts() Products {
	return productList
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (db *DB) GetProductByID(id int) (*Product, error) {
	i := db.findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (db *DB) UpdateProduct(p Product) error {
	i := db.findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &p

	return nil
}

// AddProduct adds a new product to the database
func (db *DB) AddProduct(p Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, &p)
}

// DeleteProduct deletes a product from the database
func (db *DB) DeleteProduct(id int) error {
	i := db.findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func (db *DB) findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

package data

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this product
	//
	// required: true
	// min length: 1
	Name string `json:"name" validate:"required"`

	// the description of the product
	//
	// required: false
	// max length: 1000
	Description string `json:"description"`

	// the price of the product
	//
	// required: true
	// min: 0.01
	// max: 10.00
	Price float32 `json:"price" validate:"required,gt=0"`

	// the unique stock keeping unit (SKU) for the product
	//
	// required: true
	// pattern: [a-z0-9]+\-[a-z0-9]+\-[a-z0-9]
	// example: abc-123-b4d
	SKU string `json:"sku" validate:"sku"`
}

// Products defines a slice of Product
type Products []*Product

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323-asd2-dfdf",
	},
	&Product{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34-324jf-sf12lj",
	},
}

package data

import (
	"context"
	"fmt"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

type ProductsDB struct {
	curClient    currency.CurrencyClient
	streamClient currency.Currency_SubscribeRatesClient
	log          hclog.Logger
}

// NewProductsDB returns a Data object for CRUD operations on
// Products data.
// This type also handles conversion of currencies through integraiton with the
// currency service.
func NewProductsDB(c currency.CurrencyClient, l hclog.Logger) (*ProductsDB, error) {
	sc, err := c.SubscribeRates(context.Background())
	if err != nil {
		return nil, err
	}

	pb := &ProductsDB{c, sc, l}
	pb.handleServerMessages()

	return pb, nil
}

func (p *ProductsDB) handleServerMessages() {
	go func() {
		for {
			rr, err := p.streamClient.Recv()
			if err != nil {
				p.log.Error("Received error from server", "error", err)
				break
			}

			p.log.Info("Received message from server", "Base", rr.GetBase(), "Dest", rr.GetDestination(), "Rate", rr.GetRate())
		}
	}()
}

// GetProducts returns all products from the database
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}

	return pr, nil
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	np := *productList[i]
	np.Price = np.Price * rate

	return &np, nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (p *ProductsDB) UpdateProduct(pr Product) error {
	i := findIndexByProductID(pr.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &pr

	return nil
}

// AddProduct adds a new product to the database
func (p *ProductsDB) AddProduct(pr Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	pr.ID = maxID + 1
	productList = append(productList, &pr)
}

// DeleteProduct deletes a product from the database
func (p *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	rr := &currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_value["EUR"]),
		Destination: currency.Currencies(currency.Currencies_value[destination]),
	}

	// get initial rate
	resp, err := p.curClient.GetRate(context.Background(), rr)
	if err != nil {
		return -1, fmt.Errorf("Unable to retreive exchange rate from currency service: %s", err)
	}

	// subscribe for updates
	err = p.streamClient.Send(rr)
	if err != nil {
		return -1, fmt.Errorf("Unable to subscribe for exchange rates from currency service: %s", err)
	}

	return resp.Rate, err
}

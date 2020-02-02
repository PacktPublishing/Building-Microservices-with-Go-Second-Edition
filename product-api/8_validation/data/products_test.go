package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductMissingNameReturnsErr(t *testing.T) {
	p := Product{
		Price: 1.22,
	}

	err := p.Validate()
	assert.Error(t, err)
}

func TestProductMissingPriceReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: -1,
	}

	err := p.Validate()
	assert.Error(t, err)
}

func TestProductInvalidSKUReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU: "abc",
	}

	err := p.Validate()
	assert.Error(t, err)
}

func TestValidProductDoesNOTReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU: "abc-efg-hji",
	}

	err := p.Validate()
	assert.NoError(t, err)
}
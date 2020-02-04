package data

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductMissingNameReturnsErr(t *testing.T) {
	p := Product{
		Price: 1.22,
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Error(t, err)
}

func TestProductMissingPriceReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: -1,
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Error(t, err)
}

func TestProductInvalidSKUReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU:   "abc",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Error(t, err)
}

func TestValidProductDoesNOTReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU:   "abc-efg-hji",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.NoError(t, err)
}

func TestProductsToJSON(t *testing.T) {
	ps := []*Product{
		&Product{
			Name: "abc",
		},
	}

	b := bytes.NewBufferString("")
	err := ToJSON(ps, b)
	assert.NoError(t, err)
}

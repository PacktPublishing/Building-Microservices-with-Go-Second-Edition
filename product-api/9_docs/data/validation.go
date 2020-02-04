package data

import (
	"regexp"

	"github.com/go-playground/validator"
)

// Validation contains
type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new Validation type
func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return &Validation{validate}
}

// Validate the item
// for more detail the returned error can be cast into a
// validator.ValidationErrors collection
//
// for _, errs := range err.(validator.ValidationErrors) {
//			fmt.Println(err.Namespace())
//			fmt.Println(err.Field())
//			fmt.Println(err.StructNamespace())
//			fmt.Println(err.StructField())
//			fmt.Println(err.Tag())
//			fmt.Println(err.ActualTag())
//			fmt.Println(err.Kind())
//			fmt.Println(err.Type())
//			fmt.Println(err.Value())
//			fmt.Println(err.Param())
//			fmt.Println()
//	}
func (v *Validation) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

// validateSKU
func validateSKU(fl validator.FieldLevel) bool {
	// SKU must be in the format abc-abc-abc
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	sku := re.FindAllString(fl.Field().String(), -1)

	if len(sku) == 1 {
		return true
	}

	return false
}

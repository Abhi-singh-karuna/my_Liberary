package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates struct fields
func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

// ReturnBool converts an integer to a boolean
func ReturnBool(i int) bool {
	return i == 1
}

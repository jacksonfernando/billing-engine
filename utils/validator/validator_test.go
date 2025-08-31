package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name     string  `validate:"required"`
	Age      int     `validate:"gt=0"`
	Email    string  `validate:"required,email"`
	Category string  `validate:"oneof=A B C"`
	Score    float64 `validate:"gte=0,lte=100"`
}

func TestValidateStruct_Success(t *testing.T) {
	validStruct := TestStruct{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Category: "A",
		Score:    85.5,
	}

	err := ValidateStruct(&validStruct)
	assert.NoError(t, err)
}

func TestValidateStruct_RequiredFieldMissing(t *testing.T) {
	invalidStruct := TestStruct{
		Name:     "", // Required field missing
		Age:      25,
		Email:    "john@example.com",
		Category: "A",
		Score:    85.5,
	}

	err := ValidateStruct(&invalidStruct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestValidateStruct_GreaterThanValidation(t *testing.T) {
	invalidStruct := TestStruct{
		Name:     "John Doe",
		Age:      0, // Should be greater than 0
		Email:    "john@example.com",
		Category: "A",
		Score:    85.5,
	}

	err := ValidateStruct(&invalidStruct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "age must be greater than 0")
}

func TestValidateStruct_OneOfValidation(t *testing.T) {
	invalidStruct := TestStruct{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Category: "D", // Should be one of A, B, C
		Score:    85.5,
	}

	err := ValidateStruct(&invalidStruct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category must be one of: A, B, C")
}

func TestValidateStruct_MultipleErrors(t *testing.T) {
	invalidStruct := TestStruct{
		Name:     "",              // Required field missing
		Age:      -5,              // Should be greater than 0
		Email:    "invalid-email", // Invalid email format
		Category: "D",             // Should be one of A, B, C
		Score:    150,             // Should be <= 100
	}

	err := ValidateStruct(&invalidStruct)
	assert.Error(t, err)

	errorMsg := err.Error()
	assert.Contains(t, errorMsg, "name is required")
	assert.Contains(t, errorMsg, "age must be greater than 0")
	assert.Contains(t, errorMsg, "category must be one of: A, B, C")
}

func TestValidateStruct_EmailValidation(t *testing.T) {
	invalidStruct := TestStruct{
		Name:     "John Doe",
		Age:      25,
		Email:    "invalid-email", // Invalid email format
		Category: "A",
		Score:    85.5,
	}

	err := ValidateStruct(&invalidStruct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is invalid")
}

func TestValidateStruct_RangeValidation(t *testing.T) {
	invalidStruct := TestStruct{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Category: "A",
		Score:    150, // Should be <= 100
	}

	err := ValidateStruct(&invalidStruct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "score must be less than or equal to 100")
}

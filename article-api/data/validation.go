package data

import (
	"fmt"

	"github.com/go-playground/validator"
)

// ValidationError wraps the validators FieldError
type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

// Validation contains
type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new Validation type
func NewValidation() *Validation {
	validate := validator.New()
	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	var errs validator.ValidationErrors

	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	errs = err.(validator.ValidationErrors)
	if len(errs) == 0 {
		return nil
	}

	var returnErrs []ValidationError
	for _, err := range errs {
		fmt.Printf("Validation error for param: %s", err.Param())
		// cast the FieldError into ValidationError and append to the slice
		ve := ValidationError{err}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs

}

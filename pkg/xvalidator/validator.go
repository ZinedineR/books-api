package xvalidator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"reflect"
	"regexp"
	"time"
)

// Validator is a struct that contains a pointer to a validator.Validate instance.
type Validator struct {
	validate *validator.Validate
}

// NewValidator is a function that initializes a new Validator instance.
// It registers a tag name function that returns the "name" tag of a struct field.
// It logs that the validator has been initialized and returns the new Validator instance.
func NewValidator() (*Validator, error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("name")
	})

	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		// Regular expression to check for at least one uppercase letter
		hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
		// Regular expression to check for at least one digit
		hasDigit := regexp.MustCompile(`\d`).MatchString(password)
		// Regular expression to check for at least one symbol (non-alphanumeric)
		hasSymbol := regexp.MustCompile(`[\W_]`).MatchString(password)

		// Validate the password
		return hasUpperCase && hasDigit && hasSymbol
	})

	validate.RegisterValidation("dateLocal", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()

		_, err := time.Parse("2006-01-02", dateStr)
		return err == nil

	})

	slog.Info("validator initialized")
	return &Validator{validate: validate}, nil
}

// Struct is a method of the Validator struct that validates a struct.
// It returns a slice of strings containing the validation errors.
// If there are no validation errors, it returns nil.
func (v *Validator) Struct(s interface{}) map[string]string {
	err := v.validate.Struct(s)
	if err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// Var is a method of the Validator struct that validates a single variable.
// It returns a slice of strings containing the validation errors.
// If there are no validation errors, it returns nil.
func (v *Validator) Var(field interface{}, tag string) map[string]string {
	err := v.validate.Var(field, tag)
	if err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// formatValidationError is a method of the Validator struct that formats validation errors.
// It returns a slice of strings containing the formatted validation errors.
func (v *Validator) formatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errors[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		case "email":
			errors[err.Field()] = fmt.Sprintf("%s is not a valid email", err.Field())
		case "min":
			errors[err.Field()] = fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
		case "max":
			errors[err.Field()] = fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
		case "len":
			errors[err.Field()] = fmt.Sprintf("%s must be %s characters long", err.Field(), err.Param())
		case "gte":
			errors[err.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
		case "gt":
			errors[err.Field()] = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
		case "lte":
			errors[err.Field()] = fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
		case "lt":
			errors[err.Field()] = fmt.Sprintf("%s must be less than %s", err.Field(), err.Param())
		case "numeric":
			errors[err.Field()] = fmt.Sprintf("%s must be numeric", err.Field())
		case "number":
			errors[err.Field()] = fmt.Sprintf("%s must be a number", err.Field())
		case "phone":
			errors[err.Field()] = fmt.Sprintf("%s invalid phone number", err.Field())
		case "password":
			errors[err.Field()] = fmt.Sprintf("%s is not a valid password, at least 8 characters, 1 uppercase, 1 lowercase, 1 number, and 1 special character", err.Field())
		case "dateLocal":
			errors[err.Field()] = fmt.Sprintf("%s is not a valid date, use YYYY-MM-DD", err.Field())
		default:
			errors[err.Field()] = fmt.Sprintf("%s is not valid", err.Field())
		}
	}
	return errors
}

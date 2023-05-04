package validation

import (
	"strings"
	"sync"
	"time"

	"github.com/araddon/dateparse"
	"github.com/go-playground/validator/v10"
	"github.com/thalesfsp/customerror"
)

// Singleton.
var (
	once sync.Once

	// Re-usable, cached validator.
	// SEE: https://github.com/go-playground/validator/blob/master/_examples/simple/main.go#L27
	validatorSingleton *validator.Validate
)

//////
// Built-in validators. Add more as needed.
//////

// stringContains throws an error if the field doesn't contain the string.
func stringContains(fl validator.FieldLevel) bool {
	// Get tag value.
	tagValue := fl.Param()

	// Get field value.
	fieldValue := fl.Field().String()

	// Return false if condition is met. False throws an validation error.
	return strings.Contains(fieldValue, tagValue)
}

// dateAfter returns true if the field's date is greater than or equal to the tag's date, otherwise returns false.
func dateAfter(fl validator.FieldLevel) bool {
	var finalDate time.Time

	if tagValue := fl.Param(); tagValue == "now" {
		// Subtract time to avoid errors such as (time drift).
		finalDate = time.Now().Add(time.Millisecond * -100)
	} else {
		value, err := dateparse.ParseAny(tagValue)
		if err != nil {
			return false
		}

		finalDate = value
	}

	// Get field value.
	dateField, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Return true if the field's date is greater than or equal to the tag's date.
	return dateField.Equal(finalDate) || dateField.After(finalDate)
}

// dateBefore returns true if the field's date is less than or equal to the tag's date, otherwise returns false.
func dateBefore(fl validator.FieldLevel) bool {
	var finalDate time.Time

	// Get tag value and use short syntax for the if statement.
	if tagValue := fl.Param(); tagValue == "now" {
		// Adds time to avoid errors such as (time drift).
		finalDate = time.Now().Add(time.Millisecond * 100)
	} else {
		value, err := dateparse.ParseAny(tagValue)
		if err != nil {
			return false
		}

		finalDate = value
	}

	// Get field value.
	dateField, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Return true if the field's date is less than or equal to the tag's date.
	return dateField.Equal(finalDate) || dateField.Before(finalDate)
}

//////
// Exported functionalities.
//////

// Get provides low-level access to the internal validator. Prefer to use the
// `Validate` function instead.
func Get() *validator.Validate {
	once.Do(func() {
		validatorSingleton = validator.New()

		//////
		// Register custom validators.
		//////

		if er := validatorSingleton.RegisterValidation("dateAfter", dateAfter); er != nil {
			panic("failed to register dateAfter validator")
		}

		if er := validatorSingleton.RegisterValidation("dateBefore", dateBefore); er != nil {
			panic("failed to register dateBefore validator")
		}

		if er := validatorSingleton.RegisterValidation("stringContains", stringContains); er != nil {
			panic("failed to register stringContains validator")
		}
	})

	return validatorSingleton
}

// Validate a struct.
func Validate(i any) error {
	if err := Get().Struct(i); err != nil {
		return customerror.NewInvalidError("struct", customerror.WithError(err))
	}

	return nil
}

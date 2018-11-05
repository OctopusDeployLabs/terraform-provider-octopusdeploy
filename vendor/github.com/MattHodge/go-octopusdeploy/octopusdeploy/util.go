package octopusdeploy

import (
	"fmt"
	"strings"
)

// ValidateStringInSlice checks if a string is in the given slice
func ValidateStringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}

	return false
}

// ValidatePropertyValues returns an error if the given string is not in a slice of strings
func ValidatePropertyValues(propertyName string, propertyValue string, validValues []string) error {
	if ValidateStringInSlice(propertyValue, validValues) {
		return nil
	}

	return fmt.Errorf("%s must be one of \"%v\"", propertyName, strings.Join(validValues, ","))
}

// ValidateRequiredPropertyValue returns an error if the property value is empty
func ValidateRequiredPropertyValue(propertyName string, propertyValue string) error {
	if len(propertyValue) > 0 {
		return nil
	}

	return fmt.Errorf("%s is a required property and cannot be empty", propertyName)
}

// ValidateMultipleProperties returns the first error in a list of property validations
func ValidateMultipleProperties(validatePropertyErrors []error) error {
	for _, check := range validatePropertyErrors {
		if check != nil {
			return check
		}
	}

	return nil
}

// ValidatePropertiesMatch checks two values against each other
func ValidatePropertiesMatch(firstProperty, firstPropertyName, secondProperty, secondPropertyName string) error {
	if firstProperty != secondProperty {
		return fmt.Errorf("%s and %s must match. They are currently %s and %s", firstPropertyName, secondPropertyName, firstProperty, secondProperty)
	}

	return nil
}

package domain

import (
	"database/sql/driver"
	"errors"
	"strings"
)

// ApplicationType is an enum that represents the type of application : Load balanced, Single instance, etc.
type ApplicationType string

const (
	// SingleInstance is an application type that represents an application that is not load balanced
	SingleInstance ApplicationType = "SINGLE_INSTANCE"
	// LoadBalanced is an application type that represents an application that is load balanced
	LoadBalanced ApplicationType = "LOAD_BALANCED"
)

// Scan converts the database value to the custom type
func (a *ApplicationType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("failed to scan ApplicationType: value is nil")
	}

	// Assuming the value from the database is stored as a string
	value = strings.ToUpper(string(value.([]uint8)))
	if stringValue, ok := value.(string); ok {
		*a = ApplicationType(stringValue)
		return nil
	}

	return errors.New("failed to scan ApplicationType: invalid value type")
}

// Value converts the custom type to a database value
func (a *ApplicationType) Value() (driver.Value, error) {
	return string(*a), nil
}

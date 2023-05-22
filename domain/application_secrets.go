package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// ApplicationSecret is a struct that represents a secret
type ApplicationSecret struct {
	Name string `json:"name" validate:"required"`
	Val  string `json:"value" validate:"required"`
}

// ApplicationSecrets is a slice of ApplicationSecret
type ApplicationSecrets []ApplicationSecret

func (applicationSecrets *ApplicationSecrets) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, &applicationSecrets)
}

func (applicationSecrets *ApplicationSecrets) Value() (driver.Value, error) {
	return json.Marshal(applicationSecrets)
}

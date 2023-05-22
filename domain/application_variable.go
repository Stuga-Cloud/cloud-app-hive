package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// ApplicationEnvironmentVariable is a struct that represents an environment variable
type ApplicationEnvironmentVariable struct {
	Name string `json:"name" validate:"required"`
	Val  string `json:"value" validate:"required"`
}

// ApplicationEnvironmentVariables is a slice of ApplicationEnvironmentVariable
// swagger:model ApplicationEnvironmentVariables
type ApplicationEnvironmentVariables []ApplicationEnvironmentVariable

func (applicationEnvVariables *ApplicationEnvironmentVariables) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, &applicationEnvVariables)
}

func (applicationEnvVariables *ApplicationEnvironmentVariables) Value() (driver.Value, error) {
	return json.Marshal(applicationEnvVariables)
}

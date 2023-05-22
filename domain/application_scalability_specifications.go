package domain

import (
	"cloud-app-hive/domain/errors"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// ApplicationScalabilitySpecifications is a struct that represents the scalability specifications of an application
type ApplicationScalabilitySpecifications struct {
	MinimumInstanceCount int32 `json:"minimum_instance_count" binding:"required"`
	MaximumInstanceCount int32 `json:"maximum_instance_count" binding:"required"`
	Replicas             int32 `json:"replicas" binding:"required"`
	IsAutoScaled         bool  `json:"is_auto_scaled" binding:"boolean" gorm:"default:false"`
	// If true, the application will be scaled automatically, otherwise, the user will have to scale it manually (he will be emailed -> TODO)
}

func (applicationScalabilitySpecifications ApplicationScalabilitySpecifications) Validate() error {
	if applicationScalabilitySpecifications.MinimumInstanceCount < 0 {
		return errors.NewInvalidApplicationScalabilitySpecificationsError("MinimumInstanceCount must be greater than or equal to 0")
	}
	if applicationScalabilitySpecifications.MaximumInstanceCount < 0 {
		return errors.NewInvalidApplicationScalabilitySpecificationsError("MaximumInstanceCount must be greater than or equal to 0")
	}
	if applicationScalabilitySpecifications.Replicas < 0 {
		return errors.NewInvalidApplicationScalabilitySpecificationsError("Replicas must be greater than or equal to 0")
	}
	if applicationScalabilitySpecifications.MinimumInstanceCount > applicationScalabilitySpecifications.MaximumInstanceCount {
		return errors.NewInvalidApplicationScalabilitySpecificationsError("MinimumInstanceCount must be less than or equal to MaximumInstanceCount")
	}
	return nil
}

func (applicationScalabilitySpecifications ApplicationScalabilitySpecifications) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, &applicationScalabilitySpecifications)
}

func (applicationScalabilitySpecifications ApplicationScalabilitySpecifications) Value() (driver.Value, error) {
	return json.Marshal(applicationScalabilitySpecifications)
}

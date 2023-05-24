package domain

import (
	"database/sql/driver"
	"encoding/json"

	"cloud-app-hive/controllers/errors"
)

// ApplicationScalabilitySpecifications is a struct that represents the scalability specifications of an application
// swagger:model ApplicationScalabilitySpecifications
type ApplicationScalabilitySpecifications struct {
	MinimumInstanceCount int32 `json:"minimumInstanceCount" binding:"required"`
	MaximumInstanceCount int32 `json:"maximumInstanceCount" binding:"required"`
	Replicas             int32 `json:"replicas" binding:"required"`
	IsAutoScaled         bool  `json:"isAutoScaled" binding:"boolean" gorm:"default:false"`
	// If true, the application will be scaled automatically, otherwise, the user will have to scale it manually (he will be emailed -> TODO).
}

func (applicationScalabilitySpecifications ApplicationScalabilitySpecifications) Validate() error {
	if applicationScalabilitySpecifications.MinimumInstanceCount < 0 {
		return errors.NewInvalidApplicationScalabilitySpecificationsError(
			"MinimumInstanceCount must be greater than or equal to 0",
		)
	}

	if applicationScalabilitySpecifications.MaximumInstanceCount < 0 {
		return errors.NewInvalidApplicationScalabilitySpecificationsError(
			"MaximumInstanceCount must be greater than or equal to 0",
		)
	}

	if applicationScalabilitySpecifications.Replicas < 0 {
		return errors.NewInvalidApplicationScalabilitySpecificationsError("Replicas must be greater than or equal to 0")
	}

	if applicationScalabilitySpecifications.MinimumInstanceCount > applicationScalabilitySpecifications.MaximumInstanceCount {
		return errors.NewInvalidApplicationScalabilitySpecificationsError(
			"MinimumInstanceCount must be less than or equal to MaximumInstanceCount",
		)
	}
	return nil
}

func (applicationScalabilitySpecifications ApplicationScalabilitySpecifications) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.NewInvalidApplicationScalabilitySpecificationsError(
			"failed to unmarshal JSONB value",
		)
	}

	err := json.Unmarshal(bytes, &applicationScalabilitySpecifications)
	if err != nil {
		return errors.NewInvalidApplicationScalabilitySpecificationsError(
			"failed to unmarshal JSONB value",
		)
	}

	return nil
}

func (applicationScalabilitySpecifications ApplicationScalabilitySpecifications) Value() (driver.Value, error) {
	return json.Marshal(applicationScalabilitySpecifications)
}

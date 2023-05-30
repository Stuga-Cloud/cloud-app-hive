package domain

import (
	customErrors "cloud-app-hive/controllers/errors"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// ContainerLimitUnit is an enum that represents the unit of a memory or cpu limit (e.g. KB, MB, GB, etc.)
type ContainerLimitUnit string

const (
	// KB is a limit unit that represents a limit in kilobytes.
	KB ContainerLimitUnit = "KB"
	// MB is a limit unit that represents a limit in megabytes.
	MB ContainerLimitUnit = "MB"
	// GB is a limit unit that represents a limit in gigabytes.
	GB ContainerLimitUnit = "GB"
	// TB is a limit unit that represents a limit in terabytes.
	TB ContainerLimitUnit = "TB"
	// m is a limit unit that represents a limit in millicpu.
	// The smallest allowed unit is a millicpu (m), specified as 'm'.
	// For example, '100m' is equivalent to 0.1 of a single CPU. The 'm' stands for milliCPU units.
	// So, you can specify the CPU limit as '500m', '1000m' (equivalent to 1 CPU), '2000m' (equivalent to 2 CPUs) and so on.
	m ContainerLimitUnit = "m"
)

// Scan converts the database value to the custom type
func (limitUnit *ContainerLimitUnit) Scan(value interface{}) error {
	if value == nil {
		return customErrors.NewLimitUnitScanError("failed to scan ContainerLimitUnit: value is nil")
	}

	// Assuming the value from the database is stored as a string
	if stringValue, ok := value.(string); ok {
		*limitUnit = ContainerLimitUnit(stringValue)
		return nil
	}

	return customErrors.NewLimitUnitScanError(fmt.Sprintf("failed to scan ContainerLimitUnit: value is not a string: %v", value))
}

// Value converts the custom type to a database value.
func (limitUnit *ContainerLimitUnit) Value() (driver.Value, error) {
	return string(*limitUnit), nil
}

// ContainerLimit is a struct that represents a limit of a container (e.g. CPU, Memory, Storage, etc.).
type ContainerLimit struct {
	Val  int                `json:"value" binding:"required" gorm:"not null"`
	Unit ContainerLimitUnit `json:"unit" binding:"required,oneof=KB MB GB TB m" gorm:"not null"`
}

func (containerLimit *ContainerLimit) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return customErrors.NewContainerLimitScanError(fmt.Sprintf("failed to unmarshal JSONB value: %v", value))
	}

	err := json.Unmarshal(bytes, &containerLimit)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", err)
	}

	return nil
}

func (containerLimit *ContainerLimit) Value() (driver.Value, error) {
	res, err := json.Marshal(containerLimit)
	if err != nil {
		return nil, customErrors.NewContainerLimitValueError(fmt.Sprintf("failed to marshal JSONB value: %v", err))
	}
	return res, nil
}

// ApplicationContainerSpecifications is a struct that represents the container characteristics of an application
// swagger:model ApplicationContainerSpecifications
type ApplicationContainerSpecifications struct {
	CPULimit              *ContainerLimit `json:"cpuLimit" binding:"required" gorm:"json"`
	MemoryLimit           *ContainerLimit `json:"memoryLimit" binding:"required" gorm:"json"`
	EphemeralStorageLimit *ContainerLimit `json:"ephemeralStorageLimit" binding:"required" gorm:"json"`
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return customErrors.NewApplicationContainerSpecificationsScanError(fmt.Sprintf("failed to unmarshal JSONB value: %v", value))
	}

	err := json.Unmarshal(bytes, &applicationContainerSpecifications)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", err)
	}

	return nil
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Value() (driver.Value, error) {
	res, err := json.Marshal(applicationContainerSpecifications)
	if err != nil {
		return nil, customErrors.NewApplicationContainerSpecificationsValueError(fmt.Sprintf("failed to marshal JSONB value: %v", err))
	}

	return res, nil
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Validate() error {
	if applicationContainerSpecifications.CPULimit != nil && applicationContainerSpecifications.CPULimit.Val <= 0 {
		return customErrors.NewInvalidApplicationContainerSpecificationsError("CpuLimit must be greater than 0")
	}
	if applicationContainerSpecifications.MemoryLimit != nil && applicationContainerSpecifications.MemoryLimit.Val <= 0 {
		return customErrors.NewInvalidApplicationContainerSpecificationsError("MemoryLimit must be greater than 0")
	}
	if applicationContainerSpecifications.EphemeralStorageLimit != nil && applicationContainerSpecifications.EphemeralStorageLimit.Val <= 0 {
		return customErrors.NewInvalidApplicationContainerSpecificationsError("EphemeralStorageLimit must be greater than 0")
	}
	return nil
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) SetDefaultValues() {
	if applicationContainerSpecifications.CPULimit == nil {
		applicationContainerSpecifications.CPULimit = &ContainerLimit{
			Val:  128,
			Unit: m,
		}
	}
	if applicationContainerSpecifications.MemoryLimit == nil {
		applicationContainerSpecifications.MemoryLimit = &ContainerLimit{
			Val:  128,
			Unit: MB,
		}
	}
	if applicationContainerSpecifications.EphemeralStorageLimit == nil {
		applicationContainerSpecifications.EphemeralStorageLimit = &ContainerLimit{
			Val:  256,
			Unit: MB,
		}
	}
}

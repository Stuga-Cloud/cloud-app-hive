package domain

import (
	customErrors "cloud-app-hive/controllers/errors"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// ContainerMemoryLimitUnit is an enum that represents the unit of a memory or cpu limit (e.g. KB, MB, GB, etc.)
type ContainerMemoryLimitUnit string

const (
	// KB is a limit unit that represents a limit in kilobytes.
	KB ContainerMemoryLimitUnit = "KB"
	// MB is a limit unit that represents a limit in megabytes.
	MB ContainerMemoryLimitUnit = "MB"
	// GB is a limit unit that represents a limit in gigabytes.
	GB ContainerMemoryLimitUnit = "GB"
	// TB is a limit unit that represents a limit in terabytes.
	TB ContainerMemoryLimitUnit = "TB"
)

// Scan converts the database value to the custom type
func (limitUnit *ContainerMemoryLimitUnit) Scan(value interface{}) error {
	if value == nil {
		return customErrors.NewLimitUnitScanError("failed to scan ContainerMemoryLimitUnit: value is nil")
	}

	// Assuming the value from the database is stored as a string
	if stringValue, ok := value.(string); ok {
		*limitUnit = ContainerMemoryLimitUnit(stringValue)
		return nil
	}

	return customErrors.NewLimitUnitScanError(fmt.Sprintf("failed to scan ContainerMemoryLimitUnit: value is not a string: %v", value))
}

// Value converts the custom type to a database value.
func (limitUnit *ContainerMemoryLimitUnit) Value() (driver.Value, error) {
	return string(*limitUnit), nil
}

// ContainerMemoryLimit is a struct that represents a limit of a container (e.g. CPU, Memory, Storage, etc.).
type ContainerMemoryLimit struct {
	Val  int                      `json:"value" binding:"required" gorm:"not null"`
	Unit ContainerMemoryLimitUnit `json:"unit" binding:"required,oneof=KB MB GB TB" gorm:"not null"`
}

func (containerLimit *ContainerMemoryLimit) Scan(value interface{}) error {
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

func (containerLimit *ContainerMemoryLimit) Value() (driver.Value, error) {
	res, err := json.Marshal(containerLimit)
	if err != nil {
		return nil, customErrors.NewContainerLimitValueError(fmt.Sprintf("failed to marshal JSONB value: %v", err))
	}
	return res, nil
}

type ContainerCpuLimitUnit string

const (
	// mCPU is a limit unit that represents a limit in millicpu.
	// The smallest allowed unit is a millicpu (m), specified as 'm'.
	// For example, '100m' is equivalent to 0.1 of a single CPU. The 'm' stands for milliCPU units.
	// So, you can specify the CPU limit as '500m', '1000m' (equivalent to 1 CPU), '2000m' (equivalent to 2 CPUs) and so on.
	mCPU ContainerCpuLimitUnit = "mCPU"
)

// Scan converts the database value to the custom type
func (limitUnit *ContainerCpuLimitUnit) Scan(value interface{}) error {
	if value == nil {
		return customErrors.NewLimitUnitScanError("failed to scan ContainerCpuLimitUnit: value is nil")
	}

	// Assuming the value from the database is stored as a string
	if stringValue, ok := value.(string); ok {
		*limitUnit = ContainerCpuLimitUnit(stringValue)
		return nil
	}

	return customErrors.NewLimitUnitScanError(fmt.Sprintf("failed to scan ContainerCpuLimitUnit: value is not a string: %v", value))
}

// Value converts the custom type to a database value.
func (limitUnit *ContainerCpuLimitUnit) Value() (driver.Value, error) {
	return string(*limitUnit), nil
}

type ContainerCpuLimit struct {
	Val  int                   `json:"value" binding:"required" gorm:"not null"`
	Unit ContainerCpuLimitUnit `json:"unit" binding:"required,oneof=mCPU" gorm:"not null"`
}

func (containerLimit *ContainerCpuLimit) Scan(value interface{}) error {
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

func (containerLimit *ContainerCpuLimit) Value() (driver.Value, error) {
	res, err := json.Marshal(containerLimit)
	if err != nil {
		return nil, customErrors.NewContainerLimitValueError(fmt.Sprintf("failed to marshal JSONB value: %v", err))
	}
	return res, nil
}

// ApplicationContainerSpecifications is a struct that represents the container characteristics of an application
// swagger:model ApplicationContainerSpecifications
type ApplicationContainerSpecifications struct {
	CPULimit    *ContainerCpuLimit    `json:"cpuLimit" binding:"required" gorm:"json"`
	MemoryLimit *ContainerMemoryLimit `json:"memoryLimit" binding:"required" gorm:"json"`
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
		return customErrors.NewInvalidApplicationContainerSpecificationsError("ContainerCpuLimit must be greater than 0")
	}
	if applicationContainerSpecifications.MemoryLimit != nil && applicationContainerSpecifications.MemoryLimit.Val <= 0 {
		return customErrors.NewInvalidApplicationContainerSpecificationsError("MemoryLimit must be greater than 0")
	}

	// Verify that CPU limit is contained in available choices
	if applicationContainerSpecifications.CPULimit != nil {
		isCPULimitContainedInAvailableChoices := false
		for _, cpuLimitChoice := range ApplicationCPULimitChoice {
			if applicationContainerSpecifications.CPULimit.Val == cpuLimitChoice.Value && applicationContainerSpecifications.CPULimit.Unit == cpuLimitChoice.Unit {
				isCPULimitContainedInAvailableChoices = true
				break
			}
		}
		if !isCPULimitContainedInAvailableChoices {
			return customErrors.NewInvalidApplicationContainerSpecificationsError(
				fmt.Sprintf(
					"CPU limit is not inclued in available choices: %v",
					ApplicationCPULimitChoice,
				),
			)
		}
	}

	// Verify that memory limit is contained in available choices
	if applicationContainerSpecifications.MemoryLimit != nil {
		isMemoryLimitContainedInAvailableChoices := false
		for _, memoryLimitChoice := range ApplicationMemoryLimitChoice {
			if applicationContainerSpecifications.MemoryLimit.Val == memoryLimitChoice.Value && applicationContainerSpecifications.MemoryLimit.Unit == memoryLimitChoice.Unit {
				isMemoryLimitContainedInAvailableChoices = true
				break
			}
		}
		if !isMemoryLimitContainedInAvailableChoices {
			return customErrors.NewInvalidApplicationContainerSpecificationsError(
				fmt.Sprintf(
					"Memory limit is not inclued in available choices: %v",
					ApplicationMemoryLimitChoice,
				),
			)
		}
	}

	return nil
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) SetDefaultValues() {
	if applicationContainerSpecifications.CPULimit == nil {
		applicationContainerSpecifications.CPULimit = &ContainerCpuLimit{
			Val:  50,
			Unit: mCPU,
		}
	}
	if applicationContainerSpecifications.MemoryLimit == nil {
		applicationContainerSpecifications.MemoryLimit = &ContainerMemoryLimit{
			Val:  128,
			Unit: MB,
		}
	}
}

type ApplicationCPULimit struct {
	Value int
	Unit  ContainerCpuLimitUnit
}

// Define a constant for available choices for CPU limits
var ApplicationCPULimitChoice = []ApplicationCPULimit{
	{Value: 70, Unit: mCPU},
	{Value: 140, Unit: mCPU},
	{Value: 280, Unit: mCPU},
	{Value: 560, Unit: mCPU},
	// {Value: 1120, Unit: mCPU},
	// {Value: 1680, Unit: mCPU},
	// {Value: 2240, Unit: mCPU},
}

type ApplicationMemoryLimit struct {
	Value int
	Unit  ContainerMemoryLimitUnit
}

// Define a constant for available choices for memory limits
var ApplicationMemoryLimitChoice = []ApplicationMemoryLimit{
	{Value: 128, Unit: MB},
	{Value: 256, Unit: MB},
	{Value: 512, Unit: MB},
	{Value: 1024, Unit: MB},
	// {Value: 2048, Unit: MB},
	// {Value: 4096, Unit: MB},
	// {Value: 8192, Unit: MB},
}

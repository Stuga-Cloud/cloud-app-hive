package domain

// LimitUnit is an enum that represents the unit of a limit (e.g. KB, MB, GB, etc.)
type LimitUnit int

const (
	// KB is a limit unit that represents a limit in kilobytes
	KB LimitUnit = iota
	// MB is a limit unit that represents a limit in megabytes
	MB
	// GB is a limit unit that represents a limit in gigabytes
	GB
	// TB is a limit unit that represents a limit in terabytes
	TB
)

func (limitUnit LimitUnit) String() string {
	return []string{"KB", "MB", "GB", "TB"}[limitUnit]
}

// ContainerLimit is a struct that represents a limit of a container (e.g. CPU, Memory, Storage, etc.)
type ContainerLimit struct {
	Value int
	Unit  LimitUnit
}

// ApplicationContainerSpecifications is a struct that represents the container characteristics of an application
type ApplicationContainerSpecifications struct {
	CpuLimit     ContainerLimit
	MemoryLimit  ContainerLimit
	StorageLimit ContainerLimit
}

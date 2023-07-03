package domain

type NodeMetrics struct {
	Name                          string `json:"name"`
	CPUUsage                      string `json:"cpuUsage"`
	MemoryUsage                   string `json:"memoryUsage"`
	StorageUsage                  string `json:"storageUsage"`
	EphemeralStorageUsage         string `json:"ephemeralStorageUsage"`
	Pods                          string `json:"pods"`
	ReadableCPUUsage              string `json:"readableCpuUsage"`
	ReadableMemoryUsage           string `json:"readableMemoryUsage"`
	ReadableStorageUsage          string `json:"readableStorageUsage"`
	ReadableEphemeralStorageUsage string `json:"readableEphemeralStorageUsage"`
}

type NodeCapacities struct {
	Name                     string `json:"name"`
	CPULimit                 string `json:"cpuLimit"`
	MemoryLimit              string `json:"memoryLimit"`
	StorageLimit             string `json:"storageLimit"`
	EphemeralStorageLimit    string `json:"ephemeralStorageLimit"`
	ReadableCPU              string `json:"readableCpu"`
	ReadableMemory           string `json:"readableMemory"`
	ReadableStorage          string `json:"readableStorage"`
	ReadableEphemeralStorage string `json:"readableEphemeralStorage"`
}

type NodeComputedUsage struct {
	Name                              string  `json:"name"`
	CPUUsageInPercentage              float64 `json:"cpuUsageInPercentage"`
	MemoryUsageInPercentage           float64 `json:"memoryUsageInPercentage"`
	StorageUsageInPercentage          float64 `json:"storageUsageInPercentage"`
	EphemeralStorageUsageInPercentage float64 `json:"ephemeralStorageUsageInPercentage"`
}

type ClusterMetrics struct {
	NodesMetrics        []NodeMetrics       `json:"nodesMetrics"`
	NodesCapacities     []NodeCapacities    `json:"nodesCapacities"`
	NodesComputedUsages []NodeComputedUsage `json:"computedUsage"`
}

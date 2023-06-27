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

type NodeCapacity struct {
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

type ClusterMetrics struct {
	NodesMetrics    []NodeMetrics  `json:"nodesMetrics"`
	NodesCapacities []NodeCapacity `json:"nodesCapacities"`
}

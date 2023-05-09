package domain

type ApplicationMetrics struct {
	Name                  string `json:"name"`
	CPUUsage              string `json:"cpu_usage"`
	MemoryUsage           string `json:"memory_usage"`
	StorageUsage          string `json:"storage_usage"`
	EphemeralStorageUsage string `json:"ephemeral_storage_usage"`
	PodsUsage             string `json:"pods"`
}

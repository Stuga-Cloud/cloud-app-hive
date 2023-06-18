package domain

type ApplicationLogs struct {
	PodName string `json:"podName"`
	Logs    string `json:"logs"`
}

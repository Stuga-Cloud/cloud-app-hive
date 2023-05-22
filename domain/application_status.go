package domain

type DeploymentCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
	LastUpdateTime     string `json:"last_update_time"`
	LastTransitionTime string `json:"last_transition_time"`
}

type ServiceStatus struct {
	Name string `json:"name"`
	Type string `json:"type"`
	IP   string `json:"ip"`
	Port int32  `json:"port"`
}

type IngressStatus struct {
	Name string `json:"name"`
	Type string `json:"type"`
	IP   string `json:"ip"`
	Port int32  `json:"port"`
}

type ApplicationStatus struct {
	DeploymentName      string `json:"name"`
	StatusInString      string `json:"status_in_string"`
	Replicas            int32  `json:"replicas"`
	AvailableReplicas   int32  `json:"available_replicas"`
	UnavailableReplicas int32  `json:"unavailable_replicas"`
	ReadyReplicas       int32  `json:"ready_replicas"`
	DesiredReplicas     int32  `json:"desired_replicas"`
	CurrentReplicas     int32  `json:"current_replicas"`
	UpdatedReplicas     int32  `json:"updated_replicas"`
	DeploymentCondition []DeploymentCondition
	ServiceStatus       ServiceStatus
	IngressStatus       IngressStatus
}

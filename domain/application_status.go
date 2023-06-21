package domain

import (
	"fmt"
)

type ApplicationDeploymentStatus string

const (
	AVAILABLE        ApplicationDeploymentStatus = "AVAILABLE"
	MISSING_REPLICAS ApplicationDeploymentStatus = "MISSING_REPLICAS"
	PROGRESSING      ApplicationDeploymentStatus = "PROGRESSING"
	FAILED           ApplicationDeploymentStatus = "FAILED"
	NOT_READY        ApplicationDeploymentStatus = "NOT_READY"
	UPDATING         ApplicationDeploymentStatus = "UPDATING"
)

type ContainerDeploymentStatus struct {
	Name                string                `json:"name"`
	StatusInString      string                `json:"statusInString"`
	Replicas            int                   `json:"replicas"`
	AvailableReplicas   int                   `json:"availableReplicas"`
	UnavailableReplicas int                   `json:"unavailableReplicas"`
	ReadyReplicas       int                   `json:"readyReplicas"`
	DesiredReplicas     int                   `json:"desiredReplicas"`
	CurrentReplicas     int                   `json:"currentReplicas"`
	UpdatedReplicas     int                   `json:"updatedReplicas"`
	DeploymentCondition []DeploymentCondition `json:"deploymentCondition"`
}

type DeploymentCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
	LastUpdateTime     string `json:"lastUpdateTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
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
	Name                      string                       `json:"name"`
	StatusInString            string                       `json:"statusInString"`
	Replicas                  int32                        `json:"replicas"`
	AvailableReplicas         int32                        `json:"availableReplicas"`
	UnavailableReplicas       int32                        `json:"unavailableReplicas"`
	ReadyReplicas             int32                        `json:"readyReplicas"`
	DesiredReplicas           int32                        `json:"desiredReplicas"`
	CurrentReplicas           int32                        `json:"currentReplicas"`
	UpdatedReplicas           int32                        `json:"updatedReplicas"`
	DeploymentCondition       []DeploymentCondition        `json:"deploymentCondition"`
	ComputedApplicationStatus *ApplicationDeploymentStatus `json:"computedApplicationStatus"`
	HumanizedStatus           string                       `json:"humanizedStatus"`
	ServiceStatus             ServiceStatus                `json:"serviceStatus"`
	IngressStatus             IngressStatus                `json:"ingressStatus"`
}

func (appStatus ApplicationStatus) ComputeApplicationStatus() (*ApplicationDeploymentStatus, string, error) {
	var computedStatus = FAILED
	var humanizedStatus = "DEFAULT MESSAGE"
	if appStatus.AvailableReplicas < appStatus.DesiredReplicas {
		computedStatus = MISSING_REPLICAS
		humanizedStatus = "Application is missing replicas"
	} else if appStatus.ReadyReplicas < appStatus.DesiredReplicas {
		computedStatus = NOT_READY
		humanizedStatus = "Application is not ready"
	} else if appStatus.UpdatedReplicas < appStatus.CurrentReplicas {
		computedStatus = UPDATING
		humanizedStatus = "Application is updating..."
	} else if len(appStatus.DeploymentCondition) > 0 {
		firstCondition := appStatus.DeploymentCondition[0]
		switch firstCondition.Type {
		case "Progressing":
			computedStatus = PROGRESSING
			humanizedStatus = firstCondition.Message
		case "ReplicaFailure":
			computedStatus = FAILED
			humanizedStatus = firstCondition.Message
		case "Available":
			// Do nothing
		default:
			return nil, "", fmt.Errorf("unknown deployment type: %v", firstCondition.Type)
		}
	}
	if len(appStatus.DeploymentCondition) > 0 {
		firstCondition := appStatus.DeploymentCondition[0]
		switch appStatus.DeploymentCondition[0].Type {
		case "Available":
			if firstCondition.Status == "False" {
				computedStatus = FAILED
				humanizedStatus = firstCondition.Message
			} else {
				computedStatus = AVAILABLE
				humanizedStatus = firstCondition.Message
			}
		}
	}
	
	appStatus.ComputedApplicationStatus = &computedStatus
	appStatus.HumanizedStatus = humanizedStatus

	return &computedStatus, humanizedStatus, nil
}

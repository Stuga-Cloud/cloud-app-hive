package domain

import (
	"fmt"
	"gorm.io/gorm"
	v1 "k8s.io/api/apps/v1"
	"time"
)

type Status string

const (
	AVAILABLE       Status = "AVAILABLE"
	PROGRESSING     Status = "PROGRESSING"
	REPLICA_FAILURE Status = "REPLICA_FAILURE"
	//Running  Status = "RUNNING"
	//Starting Status = "STARTING"
	//Stopping Status = "STOPPING"
	//Stopped  Status = "STOPPED"
	//Unknown  Status = "UNKNOWN"
	//Failed   Status = "FAILED"
)

func KubernetesDeploymentTypeToStatus(deploymentType string) (*Status, error) {
	var status Status
	switch deploymentType {
	case string(v1.DeploymentAvailable):
		status = AVAILABLE
	case string(v1.DeploymentProgressing):
		status = PROGRESSING
	case string(v1.DeploymentReplicaFailure):
		status = REPLICA_FAILURE
	default:
		return nil, fmt.Errorf("unknown deployment type: %v", deploymentType)
	}
	return &status, nil
}

// Application is a struct that represents a user's application
type Application struct {
	ID                        string                                `json:"id" gorm:"primaryKey"`
	Name                      string                                `json:"name" gorm:"size:100;not null"`
	Description               string                                `json:"description" gorm:"size:1000;not null"`
	Image                     string                                `json:"image" gorm:"size:1000;not null"`
	UserID                    string                                `json:"userId" gorm:"size:100;not null"`
	NamespaceID               string                                `json:"namespaceId" gorm:"size:100;not null"`
	Namespace                 Namespace                             `json:"namespace" gorm:"foreignKey:NamespaceID;references:ID;not null"`
	Port                      uint32                                `json:"port" gorm:"default:80;not null"`
	Zone                      string                                `json:"zone" gorm:"size:1000"` // The zone where the application is deployed (e.g. eu-west-1)
	ApplicationType           ApplicationType                       `json:"applicationType" gorm:"type:enum('SINGLE_INSTANCE', 'LOAD_BALANCED');default:'SINGLE_INSTANCE'"`
	EnvironmentVariables      *ApplicationEnvironmentVariables      `json:"environmentVariables" gorm:"type:json"`
	Secrets                   *ApplicationSecrets                   `json:"secrets" gorm:"type:json"`
	ContainerSpecifications   *ApplicationContainerSpecifications   `json:"containerSpecifications" gorm:"type:json"`
	ScalabilitySpecifications *ApplicationScalabilitySpecifications `json:"scalabilitySpecifications" gorm:"type:json"`
	AdministratorEmail        string                                `json:"administratorEmail" gorm:"size:320;not null"`
	Status                    *Status                               `json:"status"`
	UpdatedAt                 time.Time                             `json:"updatedAt" gorm:"autoUpdateTime;not null"`
	CreatedAt                 time.Time                             `json:"createdAt" gorm:"autoCreateTime;not null"`
	DeletedAt                 *gorm.DeletedAt                       `json:"deletedAt" gorm:"index;default:null"`
}

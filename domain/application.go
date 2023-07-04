package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ImageRegistry string

const (
	DockerHubRegistry ImageRegistry = "dockerhub"
	PrivateRegistry   ImageRegistry = "pcr"
)

// Application is a struct that represents a user's application
type Application struct {
	ID                        string                                                    `json:"id" gorm:"primaryKey"`
	Name                      string                                                    `json:"name" gorm:"size:100;not null"`
	Description               string                                                    `json:"description" gorm:"size:1000;not null"`
	Image                     string                                                    `json:"image" gorm:"size:1000;not null"`
	Registry                  ImageRegistry                                             `json:"registry" gorm:"type:enum('dockerhub', 'pcr');default:'dockerhub'"`
	UserID                    string                                                    `json:"userId" gorm:"size:100;not null"`
	NamespaceID               string                                                    `json:"namespaceId" gorm:"size:100;not null"`
	Namespace                 Namespace                                                 `json:"namespace" gorm:"foreignKey:NamespaceID;references:ID;not null"`
	Port                      uint32                                                    `json:"port" gorm:"default:80;not null"`
	Zone                      string                                                    `json:"zone" gorm:"size:1000"` // The zone where the application is deployed (e.g. eu-west-1)
	ApplicationType           ApplicationType                                           `json:"applicationType" gorm:"type:enum('SINGLE_INSTANCE', 'LOAD_BALANCED');default:'SINGLE_INSTANCE'"`
	EnvironmentVariables      *ApplicationEnvironmentVariables                          `json:"environmentVariables" gorm:"type:json"`
	Secrets                   *ApplicationSecrets                                       `json:"secrets" gorm:"type:json"`
	ContainerSpecifications   *ApplicationContainerSpecifications                       `json:"containerSpecifications" gorm:"type:json"`
	ScalabilitySpecifications *datatypes.JSONType[ApplicationScalabilitySpecifications] `json:"scalabilitySpecifications" gorm:"type:json"`
	AdministratorEmail        string                                                    `json:"administratorEmail" gorm:"size:320;not null"`
	Status                    *ApplicationDeploymentStatus                              `json:"status"`
	UpdatedAt                 time.Time                                                 `json:"updatedAt" gorm:"autoUpdateTime;not null"`
	CreatedAt                 time.Time                                                 `json:"createdAt" gorm:"autoCreateTime;not null"`
	DeletedAt                 *gorm.DeletedAt                                           `json:"deletedAt" gorm:"index;default:null"`
}

const MaxApplicationsByUser = 3
const MaxApplicationsByNamespace = 3

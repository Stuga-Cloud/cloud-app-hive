package domain

import "time"

// Application is a struct that represents a user's application
type Application struct {
	ID                        string                                `json:"id" gorm:"primaryKey"`
	Name                      string                                `json:"name" gorm:"uniqueIndex:idx_name_namespace_id;size:100;not null"`
	Description               string                                `json:"description" gorm:"size:1000;not null"`
	Image                     string                                `json:"image" gorm:"size:1000;not null"`
	UserID                    string                                `json:"user_id" gorm:"size:100;not null"`
	NamespaceID               string                                `json:"namespace_id" gorm:"uniqueIndex:idx_name_namespace_id;size:100;not null"`
	Namespace                 Namespace                             `json:"namespace" gorm:"foreignKey:NamespaceID;references:ID;not null"`
	Port                      uint32                                `json:"port" gorm:"default:80;not null"`
	Zone                      string                                `json:"zone" gorm:"size:1000"` // The zone where the application is deployed (e.g. eu-west-1)
	ApplicationType           ApplicationType                       `json:"application_type" gorm:"type:enum('SINGLE_INSTANCE', 'LOAD_BALANCED');default:'SINGLE_INSTANCE'"`
	EnvironmentVariables      *ApplicationEnvironmentVariables      `json:"environment_variables" gorm:"type:json"`
	Secrets                   *ApplicationSecrets                   `json:"secrets" gorm:"type:json"`
	ContainerSpecifications   *ApplicationContainerSpecifications   `json:"container_specifications" gorm:"type:json"`
	ScalabilitySpecifications *ApplicationScalabilitySpecifications `json:"scalability_specifications" gorm:"type:json"`
	UpdatedAt                 time.Time                             `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt                 time.Time                             `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt                 time.Time                             `json:"deleted_at" gorm:"index;default:null"`
}

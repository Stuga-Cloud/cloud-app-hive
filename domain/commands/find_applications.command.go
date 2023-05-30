package commands

import "cloud-app-hive/domain"

// FindApplications is a command that represents the parameters for finding applications
type FindApplications struct {
	Name            *string
	Image           *string
	NamespaceID     *string
	ApplicationType *domain.ApplicationType
	IsAutoScaled    *bool
	Page            uint32
	Limit           uint32
}

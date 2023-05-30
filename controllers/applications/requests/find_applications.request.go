package requests

import (
	"cloud-app-hive/domain"
)

// FindApplicationsRequest is a struct that represents the request body for finding applications
// swagger:model FindApplicationsRequest
type FindApplicationsRequest struct {
	Name            *string                 `form:"name" binding:""`
	Image           *string                 `form:"image" binding:""`
	NamespaceID     *string                 `form:"namespaceId" binding:""`
	ApplicationType *domain.ApplicationType `form:"applicationType" binding:"omitempty,oneof=SINGLE_INSTANCE LOAD_BALANCED"`
	IsAutoScaled    *bool                   `form:"isAutoScaled" binding:""`
	Page            uint32                  `form:"page" binding:"required"`
	Limit           uint32                  `form:"limit" binding:"required"`
}

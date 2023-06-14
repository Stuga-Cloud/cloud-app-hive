package requests

import "cloud-app-hive/domain"

// AddMemberToNamespaceRequest is a struct that represents the request body for creating a namespace
type AddMemberToNamespaceRequest struct {
	UserID  string      `json:"userId" binding:"required"`
	Role    domain.Role `json:"role" binding:"required,oneof=ADMIN MEMBER"`
	AddedBy string      `json:"addedBy" binding:"required"`
}

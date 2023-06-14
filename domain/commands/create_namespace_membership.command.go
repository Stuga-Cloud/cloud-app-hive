package commands

import "cloud-app-hive/domain"

// CreateNamespaceMembership is a command that represents the intent to create a namespace membership.
type CreateNamespaceMembership struct {
	UserID      string      `json:"user_id" binding:"required"`
	NamespaceID string      `json:"namespace_id" binding:"required"`
	Role        domain.Role `json:"role" binding:"required,oneof=ADMIN MEMBER"`
	AddedBy     string      `json:"added_by" binding:"required"`
}

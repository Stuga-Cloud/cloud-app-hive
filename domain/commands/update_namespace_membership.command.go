package commands

import "cloud-app-hive/domain"

// UpdateNamespaceMembership is a command that represents the intent to create a namespace membership.
type UpdateNamespaceMembership struct {
	UserID      string      `json:"user_id" binding:"required"`
	NamespaceID string      `json:"namespace_id" binding:"required"`
	Role        domain.Role `json:"role" binding:"required,oneof=ADMIN MEMBER"`
}

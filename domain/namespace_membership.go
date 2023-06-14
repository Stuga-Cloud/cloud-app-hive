package domain

import (
	"gorm.io/gorm"
	"time"
)

type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleMember Role = "MEMBER"
)

// NamespaceMembership is a struct that represents a user's membership to a namespace
type NamespaceMembership struct {
	ID          string          `json:"id" gorm:"primaryKey"`
	UserID      string          `json:"userId" gorm:"not null"`
	NamespaceID string          `json:"namespaceId" gorm:"size:100;not null"`
	Role        Role            `json:"role" gorm:"enum:ADMIN,MEMBER;not null"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"autoUpdateTime;not null"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"autoCreateTime;not null"`
	DeletedAt   *gorm.DeletedAt `json:"deletedAt" gorm:"index;default:null"`
}

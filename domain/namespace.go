package domain

import (
	"gorm.io/gorm"
	"time"
)

// Namespace is a struct that represents a user's namespace
type Namespace struct {
	ID           string                `json:"id" gorm:"primaryKey"`
	Name         string                `json:"name" gorm:"size:100;not null"`
	Description  string                `json:"description" gorm:"size:1000"`
	UserID       string                `json:"userId" gorm:"index:idx_user_id;not null"`
	Memberships  []NamespaceMembership `json:"memberships" gorm:"foreignKey:NamespaceID;references:ID;not null"`
	Applications []Application         `json:"applications" gorm:"foreignKey:NamespaceID;references:ID;not null"`
	UpdatedAt    time.Time             `json:"updatedAt" gorm:"autoUpdateTime"`
	CreatedAt    time.Time             `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt    *gorm.DeletedAt       `json:"deletedAt" gorm:"index;default:null"`
}

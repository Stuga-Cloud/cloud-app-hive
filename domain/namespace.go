package domain

import "time"

// Namespace is a struct that represents a user's namespace
type Namespace struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;size:100;not null"`
	Description string    `json:"description" gorm:"size:1000"`
	UserID      string    `json:"userId" gorm:"index:idx_user_id;not null"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt   time.Time `json:"deletedAt" gorm:"index;default:null"`
}

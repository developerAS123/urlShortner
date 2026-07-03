package models

import (
	"time"

	"gorm.io/gorm"
)

type Link struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	OriginalURL string         `gorm:"not null" json:"original_url"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Title     string         `binding:"required,min=2,max=32" json:"title"`
	Author    string         `binding:"required,min=8,max=32" json:"author"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID    uint           `json:"user_id"`
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Stock       int32          `gorm:"not null" json:"stock"`
	Price       float32        `gorm:"not null" json:"price"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Users       []User         `gorm:"many2many:enrollments;" json:"users"`
}

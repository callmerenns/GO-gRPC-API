package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string         `gorm:"not null;column:firstname" json:"firstname"`
	LastName  string         `gorm:"not null;column:lastname" json:"lastname"`
	Email     string         `gorm:"not null;unique" json:"email"`
	Password  string         `gorm:"not null" json:"password"`
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Products  []Product      `gorm:"many2many:enrollments;" json:"products"`
}

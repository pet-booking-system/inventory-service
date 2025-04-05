package models

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ResourceID  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"not null"`
	Type        string    `gorm:"index"`
	Status      string    `gorm:"type:resource_status;not null;default:'available'"`
	Description string
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}

package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Model is a base model for db interaction
type Model struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// BeforeCreate implements UUID generation for new resources
func (*Model) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New())
}

// service is a convenience struct for grouping methods
type service struct{}
